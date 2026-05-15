package command

import (
	"context"
	"errors"

	"vicomova/internal/chat/domain/entity"
	"vicomova/internal/chat/domain/repository"
	ws "vicomova/internal/chat/interfaces/ws"
	"vicomova/pkg/log"
)

var (
	ErrNotFollowing     = errors.New("not following")
	ErrCannotFollowSelf = errors.New("cannot follow self")
	ErrNotFriends       = errors.New("not friends")
	ErrNotGroupMember   = errors.New("not group member")
	ErrNotGroupAdmin    = errors.New("not group admin")
	ErrAlreadyFollowed  = errors.New("already followed")
)

type ChatCommandService struct {
	followRepo         repository.FollowRepository
	conversationRepo   repository.ConversationRepository
	groupMemberRepo    repository.GroupMemberRepository
	messageRepo        repository.MessageRepository
	offlineMessageRepo repository.OfflineMessageRepository
}

func NewChatCommandService(
	followRepo repository.FollowRepository,
	conversationRepo repository.ConversationRepository,
	groupMemberRepo repository.GroupMemberRepository,
	messageRepo repository.MessageRepository,
	offlineMessageRepo repository.OfflineMessageRepository,
) *ChatCommandService {
	return &ChatCommandService{
		followRepo:         followRepo,
		conversationRepo:   conversationRepo,
		groupMemberRepo:    groupMemberRepo,
		messageRepo:        messageRepo,
		offlineMessageRepo: offlineMessageRepo,
	}
}

func (s *ChatCommandService) Follow(ctx context.Context, cmd *FollowCommand) error {
	if cmd.UserID == cmd.TargetID {
		return ErrCannotFollowSelf
	}

	follow := &entity.Follow{
		FollowerID:  cmd.UserID,
		FollowingID: cmd.TargetID,
	}
	if err := s.followRepo.Create(ctx, follow); err != nil {
		log.Error.Printf("ChatCommandService.Follow: failed: %v", err)
		return err
	}
	return nil
}

func (s *ChatCommandService) Unfollow(ctx context.Context, cmd *UnfollowCommand) error {
	isFollowing, err := s.followRepo.IsFollowing(ctx, cmd.UserID, cmd.TargetID)
	if err != nil {
		log.Error.Printf("ChatCommandService.Unfollow: IsFollowing failed: %v", err)
		return err
	}
	if !isFollowing {
		return ErrNotFollowing
	}

	if err := s.followRepo.Delete(ctx, cmd.UserID, cmd.TargetID); err != nil {
		log.Error.Printf("ChatCommandService.Unfollow: Delete failed: %v", err)
		return err
	}
	return nil
}

func (s *ChatCommandService) SendMessage(ctx context.Context, cmd *SendMessageCommand) (*SendMessageResult, error) {
	isFollowing1, err := s.followRepo.IsFollowing(ctx, cmd.SenderID, cmd.ReceiverID)
	if err != nil {
		log.Error.Printf("ChatCommandService.SendMessage: IsFollowing1 failed: %v", err)
		return nil, err
	}
	isFollowing2, err := s.followRepo.IsFollowing(ctx, cmd.ReceiverID, cmd.SenderID)
	if err != nil {
		log.Error.Printf("ChatCommandService.SendMessage: IsFollowing2 failed: %v", err)
		return nil, err
	}

	if !isFollowing1 || !isFollowing2 {
		return nil, ErrNotFriends
	}

	conv, err := s.conversationRepo.Get1v1(ctx, cmd.SenderID, cmd.ReceiverID)
	if err != nil {
		log.Error.Printf("ChatCommandService.SendMessage: Get1v1 failed: %v", err)
		return nil, err
	}

	if conv == nil {
		conv, err = s.conversationRepo.Create1v1(ctx, cmd.SenderID, cmd.ReceiverID)
		if err != nil {
			log.Error.Printf("ChatCommandService.SendMessage: Create1v1 failed: %v", err)
			return nil, err
		}
	}

	msg := &entity.Message{
		ConversationID: conv.ID,
		SenderID:       cmd.SenderID,
		Content:        cmd.Content,
	}
	if err := s.messageRepo.Create(ctx, msg); err != nil {
		log.Error.Printf("ChatCommandService.SendMessage: Create failed: %v", err)
		return nil, err
	}

	// 实时推送或存储离线消息
	s.deliverOrStoreMessage(cmd.ReceiverID, cmd.SenderID, cmd.Content, msg.ID)

	return &SendMessageResult{
		MessageID:      msg.ID,
		ConversationID: conv.ID,
	}, nil
}

func (s *ChatCommandService) deliverOrStoreMessage(toUserID, fromUserID int64, content string, msgID int64) {
	// 检查用户是否在线
	if ws.IsUserOnline(toUserID) {
		// 在线，推送消息
		ws.PushMessageToUser(toUserID, fromUserID, content, msgID)
		log.Info.Printf("Message pushed to online user: toUserID=%d, fromUserID=%d, msgID=%d", toUserID, fromUserID, msgID)
	} else {
		// 离线，存储消息
		offlineMsg := &entity.OfflineMessage{
			ToUserID:   toUserID,
			FromUserID: fromUserID,
			Content:    content,
		}
		if err := s.offlineMessageRepo.Create(context.Background(), offlineMsg); err != nil {
			log.Error.Printf("Failed to store offline message: %v", err)
		} else {
			log.Info.Printf("Message stored offline: toUserID=%d, fromUserID=%d, msgID=%d", toUserID, fromUserID, msgID)
		}
	}
}

func (s *ChatCommandService) CreateGroup(ctx context.Context, cmd *CreateGroupCommand) (*CreateGroupResult, error) {
	conv := &entity.Conversation{
		Type: entity.ConversationTypeGroup,
		Name: &cmd.Name,
	}

	var err error
	if err = s.conversationRepo.Create(ctx, conv); err != nil {
		log.Error.Printf("ChatCommandService.CreateGroup: Create failed: %v", err)
		return nil, err
	}

	member := &entity.GroupMember{
		GroupID: conv.ID,
		UserID:  cmd.CreatorID,
		Role:    entity.GroupRoleAdmin,
	}
	if err = s.groupMemberRepo.Create(ctx, member); err != nil {
		log.Error.Printf("ChatCommandService.CreateGroup: Create member failed: %v", err)
		return nil, err
	}

	return &CreateGroupResult{
		GroupID:        conv.ID,
		ConversationID: conv.ID,
	}, nil
}

func (s *ChatCommandService) AddGroupMember(ctx context.Context, cmd *AddGroupMemberCommand) error {
	isAdmin, err := s.groupMemberRepo.IsAdmin(ctx, cmd.GroupID, cmd.AdminID)
	if err != nil {
		log.Error.Printf("ChatCommandService.AddGroupMember: IsAdmin failed: %v", err)
		return err
	}
	if !isAdmin {
		return ErrNotGroupAdmin
	}

	member := &entity.GroupMember{
		GroupID: cmd.GroupID,
		UserID:  cmd.UserID,
		Role:    entity.GroupRoleMember,
	}
	if err := s.groupMemberRepo.Create(ctx, member); err != nil {
		log.Error.Printf("ChatCommandService.AddGroupMember: Create failed: %v", err)
		return err
	}
	return nil
}

func (s *ChatCommandService) RemoveGroupMember(ctx context.Context, cmd *RemoveGroupMemberCommand) error {
	isAdmin, err := s.groupMemberRepo.IsAdmin(ctx, cmd.GroupID, cmd.AdminID)
	if err != nil {
		log.Error.Printf("ChatCommandService.RemoveGroupMember: IsAdmin failed: %v", err)
		return err
	}
	if !isAdmin {
		return ErrNotGroupAdmin
	}

	isMember, err := s.groupMemberRepo.IsMember(ctx, cmd.GroupID, cmd.UserID)
	if err != nil {
		log.Error.Printf("ChatCommandService.RemoveGroupMember: IsMember failed: %v", err)
		return err
	}
	if !isMember {
		return ErrNotGroupMember
	}

	if err := s.groupMemberRepo.Delete(ctx, cmd.GroupID, cmd.UserID); err != nil {
		log.Error.Printf("ChatCommandService.RemoveGroupMember: Delete failed: %v", err)
		return err
	}
	return nil
}

func (s *ChatCommandService) SendGroupMessage(ctx context.Context, cmd *SendGroupMessageCommand) (*SendGroupMessageResult, error) {
	isMember, err := s.groupMemberRepo.IsMember(ctx, cmd.GroupID, cmd.SenderID)
	if err != nil {
		log.Error.Printf("ChatCommandService.SendGroupMessage: IsMember failed: %v", err)
		return nil, err
	}
	if !isMember {
		return nil, ErrNotGroupMember
	}

	msg := &entity.Message{
		ConversationID: cmd.GroupID,
		SenderID:       cmd.SenderID,
		Content:        cmd.Content,
	}
	if err := s.messageRepo.Create(ctx, msg); err != nil {
		log.Error.Printf("ChatCommandService.SendGroupMessage: Create failed: %v", err)
		return nil, err
	}

	// 获取群成员并推送消息
	members, err := s.groupMemberRepo.ListByGroup(ctx, cmd.GroupID)
	if err != nil {
		log.Error.Printf("ChatCommandService.SendGroupMessage: ListByGroup failed: %v", err)
		return nil, err
	}

	memberIDs := make([]int64, len(members))
	for i, m := range members {
		memberIDs[i] = m.UserID
	}

	// 群消息推送
	ws.PushGroupMessage(cmd.GroupID, cmd.SenderID, cmd.Content, msg.ID, memberIDs)
	log.Info.Printf("Group message sent: groupID=%d, senderID=%d, msgID=%d, members=%d", cmd.GroupID, cmd.SenderID, msg.ID, len(memberIDs))

	return &SendGroupMessageResult{
		MessageID:      msg.ID,
		ConversationID: cmd.GroupID,
	}, nil
}
