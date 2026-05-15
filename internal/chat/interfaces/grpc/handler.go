package grpc

import (
	"context"
	"strconv"

	"vicomova/internal/chat/application/command"
	"vicomova/internal/chat/application/query"
	userRpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/pkg/log"
	chat "vicomova/third_party/kitex_gen/chat"
)

type ChatHandler struct {
	cmdSvc  *command.ChatCommandService
	qrySvc  *query.ChatQueryService
	userCli *userRpc.UserClient
}

func NewChatHandler(cmdSvc *command.ChatCommandService, qrySvc *query.ChatQueryService, userCli *userRpc.UserClient) *ChatHandler {
	return &ChatHandler{
		cmdSvc:  cmdSvc,
		qrySvc:  qrySvc,
		userCli: userCli,
	}
}

func (h *ChatHandler) Follow(ctx context.Context, req *chat.FollowRequest) (*chat.FollowResponse, error) {
	cmd := &command.FollowCommand{
		UserID:   req.UserId,
		TargetID: req.TargetId,
	}
	if err := h.cmdSvc.Follow(ctx, cmd); err != nil {
		return nil, err
	}
	return &chat.FollowResponse{}, nil
}

func (h *ChatHandler) Unfollow(ctx context.Context, req *chat.UnfollowRequest) (*chat.UnfollowResponse, error) {
	cmd := &command.UnfollowCommand{
		UserID:   req.UserId,
		TargetID: req.TargetId,
	}
	if err := h.cmdSvc.Unfollow(ctx, cmd); err != nil {
		return nil, err
	}
	return &chat.UnfollowResponse{}, nil
}

func (h *ChatHandler) ListFollowers(ctx context.Context, req *chat.ListFollowersRequest) (*chat.ListFollowersResponse, error) {
	followerIDs, err := h.qrySvc.ListFollowers(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if len(followerIDs) == 0 {
		return &chat.ListFollowersResponse{Followers: []*chat.User{}}, nil
	}

	userMap, err := h.userCli.BatchGetUsers(ctx, followerIDs)
	if err != nil {
		log.Error.Printf("ChatHandler.ListFollowers: BatchGetUsers failed: %v", err)
		return nil, err
	}

	protoFollowers := make([]*chat.User, 0, len(followerIDs))
	for _, id := range followerIDs {
		if u, ok := userMap[id]; ok {
			protoFollowers = append(protoFollowers, &chat.User{
				Id:       u.UserId,
				Username: u.Username,
				Avatar:   "",
			})
		}
	}

	return &chat.ListFollowersResponse{Followers: protoFollowers}, nil
}

func (h *ChatHandler) ListFollowing(ctx context.Context, req *chat.ListFollowingRequest) (*chat.ListFollowingResponse, error) {
	followingIDs, err := h.qrySvc.ListFollowing(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if len(followingIDs) == 0 {
		return &chat.ListFollowingResponse{Following: []*chat.User{}}, nil
	}

	userMap, err := h.userCli.BatchGetUsers(ctx, followingIDs)
	if err != nil {
		log.Error.Printf("ChatHandler.ListFollowing: BatchGetUsers failed: %v", err)
		return nil, err
	}

	protoFollowing := make([]*chat.User, 0, len(followingIDs))
	for _, id := range followingIDs {
		if u, ok := userMap[id]; ok {
			protoFollowing = append(protoFollowing, &chat.User{
				Id:       u.UserId,
				Username: u.Username,
				Avatar:   "",
			})
		}
	}

	return &chat.ListFollowingResponse{Following: protoFollowing}, nil
}

func (h *ChatHandler) ListFriends(ctx context.Context, req *chat.ListFriendsRequest) (*chat.ListFriendsResponse, error) {
	friendIDs, err := h.qrySvc.ListFriends(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if len(friendIDs) == 0 {
		return &chat.ListFriendsResponse{Friends: []*chat.User{}}, nil
	}

	userMap, err := h.userCli.BatchGetUsers(ctx, friendIDs)
	if err != nil {
		log.Error.Printf("ChatHandler.ListFriends: BatchGetUsers failed: %v", err)
		return nil, err
	}

	protoFriends := make([]*chat.User, 0, len(friendIDs))
	for _, id := range friendIDs {
		if u, ok := userMap[id]; ok {
			protoFriends = append(protoFriends, &chat.User{
				Id:       u.UserId,
				Username: u.Username,
				Avatar:   "",
			})
		}
	}

	return &chat.ListFriendsResponse{Friends: protoFriends}, nil
}

func (h *ChatHandler) IsFriends(ctx context.Context, req *chat.IsFriendsRequest) (*chat.IsFriendsResponse, error) {
	isFriend1, err := h.qrySvc.IsFollowing(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, err
	}
	isFriend2, err := h.qrySvc.IsFollowing(ctx, req.TargetId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &chat.IsFriendsResponse{IsFriends: isFriend1 && isFriend2}, nil
}

func (h *ChatHandler) IsGroupMember(ctx context.Context, req *chat.IsGroupMemberRequest) (*chat.IsGroupMemberResponse, error) {
	isMember, err := h.qrySvc.IsGroupMember(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	return &chat.IsGroupMemberResponse{IsMember: isMember}, nil
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	cmd := &command.SendMessageCommand{
		SenderID:   req.SenderId,
		ReceiverID: req.ReceiverId,
		Content:    req.Content,
	}
	result, err := h.cmdSvc.SendMessage(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &chat.SendMessageResponse{
		MessageId:      result.MessageID,
		ConversationId: result.ConversationID,
	}, nil
}

// SendMessageWS implements ws.ChatSender for WebSocket
func (h *ChatHandler) SendMessageWS(ctx context.Context, senderID, receiverID int64, content string) error {
	cmd := &command.SendMessageCommand{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	_, err := h.cmdSvc.SendMessage(ctx, cmd)
	return err
}

// SendGroupMessageWS implements ws.ChatGroupSender for WebSocket
func (h *ChatHandler) SendGroupMessageWS(ctx context.Context, senderID, groupID int64, content string) error {
	cmd := &command.SendGroupMessageCommand{
		SenderID: senderID,
		GroupID:  groupID,
		Content:  content,
	}
	_, err := h.cmdSvc.SendGroupMessage(ctx, cmd)
	return err
}

func (h *ChatHandler) ListMessages(ctx context.Context, req *chat.ListMessagesRequest) (*chat.ListMessagesResponse, error) {
	conv, err := h.qrySvc.Get1v1Conversation(ctx, req.UserId, req.PeerId)
	if err != nil {
		return nil, err
	}

	if conv == nil {
		return &chat.ListMessagesResponse{
			Messages:   []*chat.Message{},
			NextCursor: "",
			HasMore:    false,
		}, nil
	}

	messages, hasMore, err := h.qrySvc.ListMessages(ctx, conv.ID, req.Cursor, int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*chat.Message, len(messages))
	for i, m := range messages {
		protoMessages[i] = &chat.Message{
			Id:             m.ID,
			ConversationId: m.ConversationID,
			SenderId:       m.SenderID,
			Content:        m.Content,
			CreatedAt:      m.CreatedAt.Unix(),
		}
	}

	var nextCursor string
	if hasMore && len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		nextCursor = strconv.FormatInt(lastMsg.CreatedAt.UnixMilli(), 10)
	}

	return &chat.ListMessagesResponse{
		Messages:   protoMessages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (h *ChatHandler) CreateGroup(ctx context.Context, req *chat.CreateGroupRequest) (*chat.CreateGroupResponse, error) {
	cmd := &command.CreateGroupCommand{
		CreatorID: req.CreatorId,
		Name:      req.Name,
	}
	result, err := h.cmdSvc.CreateGroup(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &chat.CreateGroupResponse{
		GroupId:        result.GroupID,
		ConversationId: result.ConversationID,
	}, nil
}

func (h *ChatHandler) AddGroupMember(ctx context.Context, req *chat.AddGroupMemberRequest) (*chat.AddGroupMemberResponse, error) {
	cmd := &command.AddGroupMemberCommand{
		GroupID: req.GroupId,
		UserID:  req.UserId,
		AdminID: req.AdminId,
	}
	if err := h.cmdSvc.AddGroupMember(ctx, cmd); err != nil {
		return nil, err
	}
	return &chat.AddGroupMemberResponse{}, nil
}

func (h *ChatHandler) RemoveGroupMember(ctx context.Context, req *chat.RemoveGroupMemberRequest) (*chat.RemoveGroupMemberResponse, error) {
	cmd := &command.RemoveGroupMemberCommand{
		GroupID: req.GroupId,
		UserID:  req.UserId,
		AdminID: req.AdminId,
	}
	if err := h.cmdSvc.RemoveGroupMember(ctx, cmd); err != nil {
		return nil, err
	}
	return &chat.RemoveGroupMemberResponse{}, nil
}

func (h *ChatHandler) ListGroupMembers(ctx context.Context, req *chat.ListGroupMembersRequest) (*chat.ListGroupMembersResponse, error) {
	members, err := h.qrySvc.ListGroupMembers(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	protoMembers := make([]*chat.GroupMember, len(members))
	for i, m := range members {
		protoMembers[i] = &chat.GroupMember{
			UserId:   m.UserID,
			Role:     int64(m.Role),
			JoinedAt: m.JoinedAt.Unix(),
		}
	}

	return &chat.ListGroupMembersResponse{
		Members: protoMembers,
	}, nil
}

func (h *ChatHandler) ListUserGroups(ctx context.Context, req *chat.ListUserGroupsRequest) (*chat.ListUserGroupsResponse, error) {
	groupIDs, err := h.qrySvc.ListUserGroups(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	protoGroups := make([]*chat.Group, len(groupIDs))
	for i, id := range groupIDs {
		protoGroups[i] = &chat.Group{
			Id: id,
		}
	}

	return &chat.ListUserGroupsResponse{
		Groups: protoGroups,
	}, nil
}

func (h *ChatHandler) SendGroupMessage(ctx context.Context, req *chat.SendGroupMessageRequest) (*chat.SendGroupMessageResponse, error) {
	cmd := &command.SendGroupMessageCommand{
		SenderID: req.SenderId,
		GroupID:  req.GroupId,
		Content:  req.Content,
	}
	result, err := h.cmdSvc.SendGroupMessage(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &chat.SendGroupMessageResponse{
		MessageId:      result.MessageID,
		ConversationId: result.ConversationID,
	}, nil
}

func (h *ChatHandler) ListGroupMessages(ctx context.Context, req *chat.ListGroupMessagesRequest) (*chat.ListGroupMessagesResponse, error) {
	messages, hasMore, err := h.qrySvc.ListGroupMessages(ctx, req.GroupId, req.Cursor, int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*chat.Message, len(messages))
	for i, m := range messages {
		protoMessages[i] = &chat.Message{
			Id:             m.ID,
			ConversationId: m.ConversationID,
			SenderId:       m.SenderID,
			Content:        m.Content,
			CreatedAt:      m.CreatedAt.Unix(),
		}
	}

	var nextCursor string
	if hasMore && len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		nextCursor = strconv.FormatInt(lastMsg.CreatedAt.UnixMilli(), 10)
	}

	return &chat.ListGroupMessagesResponse{
		Messages:   protoMessages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
