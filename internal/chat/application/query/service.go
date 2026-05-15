package query

import (
	"context"

	"vicomova/internal/chat/domain/entity"
	"vicomova/internal/chat/domain/repository"
)

type ChatQueryService struct {
	followRepo         repository.FollowRepository
	groupMemberRepo    repository.GroupMemberRepository
	messageRepo        repository.MessageRepository
	offlineMessageRepo repository.OfflineMessageRepository
}

func NewChatQueryService(
	followRepo repository.FollowRepository,
	groupMemberRepo repository.GroupMemberRepository,
	messageRepo repository.MessageRepository,
	offlineMessageRepo repository.OfflineMessageRepository,
) *ChatQueryService {
	return &ChatQueryService{
		followRepo:         followRepo,
		groupMemberRepo:    groupMemberRepo,
		messageRepo:        messageRepo,
		offlineMessageRepo: offlineMessageRepo,
	}
}

func (s *ChatQueryService) ListFollowers(ctx context.Context, userID int64) ([]int64, error) {
	return s.followRepo.ListFollowers(ctx, userID)
}

func (s *ChatQueryService) ListFollowing(ctx context.Context, userID int64) ([]int64, error) {
	return s.followRepo.ListFollowing(ctx, userID)
}

func (s *ChatQueryService) ListFriends(ctx context.Context, userID int64) ([]int64, error) {
	return s.followRepo.ListFriends(ctx, userID)
}

func (s *ChatQueryService) IsFollowing(ctx context.Context, userID, targetID int64) (bool, error) {
	return s.followRepo.IsFollowing(ctx, userID, targetID)
}

func (s *ChatQueryService) ListMessages(ctx context.Context, conversationID int64, cursor int64, limit int) ([]*entity.Message, bool, error) {
	return s.messageRepo.ListByConversation(ctx, conversationID, cursor, limit)
}

func (s *ChatQueryService) ListGroupMembers(ctx context.Context, groupID int64) ([]*entity.GroupMember, error) {
	return s.groupMemberRepo.ListByGroup(ctx, groupID)
}

func (s *ChatQueryService) ListUserGroups(ctx context.Context, userID int64) ([]int64, error) {
	return s.groupMemberRepo.ListByUser(ctx, userID)
}

func (s *ChatQueryService) ListGroupMessages(ctx context.Context, groupID int64, cursor int64, limit int) ([]*entity.Message, bool, error) {
	return s.messageRepo.ListByConversation(ctx, groupID, cursor, limit)
}

func (s *ChatQueryService) IsGroupMember(ctx context.Context, userID, groupID int64) (bool, error) {
	return s.groupMemberRepo.IsMember(ctx, groupID, userID)
}

func (s *ChatQueryService) Get1v1Conversation(ctx context.Context, userID1, userID2 int64) (*entity.Conversation, error) {
	return nil, nil
}

// GetOfflineMessages 获取离线消息
func (s *ChatQueryService) GetOfflineMessages(ctx context.Context, userID int64, limit int) ([]*entity.OfflineMessage, error) {
	return s.offlineMessageRepo.ListByUser(ctx, userID, limit)
}

// MarkOfflineMessagesAsRead 标记离线消息已读
func (s *ChatQueryService) MarkOfflineMessagesAsRead(ctx context.Context, userID int64) error {
	return s.offlineMessageRepo.MarkAsRead(ctx, userID)
}

// GetUnreadCount 获取未读消息数
func (s *ChatQueryService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.offlineMessageRepo.UnreadCount(ctx, userID)
}
