package handler

import (
	"context"
	"fmt"
	"strconv"

	chatRpc "vicomova/internal/chat/interfaces/grpc"
	ws "vicomova/internal/chat/interfaces/ws"
	"vicomova/pkg/constants"
	hertz "vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type ChatHandler struct {
	chatClient *chatRpc.ChatClient
}

func NewChatHandler(chatClient *chatRpc.ChatClient) *ChatHandler {
	return &ChatHandler{chatClient: chatClient}
}

// ========== 关注相关 ==========

// @Summary 关注用户
// @Description 关注另一个用户（需要认证）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body FollowRequest true "关注请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/follow [post]
func (h *ChatHandler) Follow(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req FollowRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.chatClient.Follow(ctx, userID, req.TargetID)
	if err != nil {
		hlog.Errorf("Follow: userID=%d targetID=%d failed: %v", userID, req.TargetID, err)
		c.JSON(500, hertz.Fail(500, "Failed to follow"))
		return
	}

	hlog.Infof("Follow: userID=%d targetID=%d success", userID, req.TargetID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 取消关注
// @Description 取消关注用户（需要认证）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UnfollowRequest true "取消关注请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/follow [delete]
func (h *ChatHandler) Unfollow(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req UnfollowRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.chatClient.Unfollow(ctx, userID, req.TargetID)
	if err != nil {
		hlog.Errorf("Unfollow: userID=%d targetID=%d failed: %v", userID, req.TargetID, err)
		c.JSON(500, hertz.Fail(500, "Failed to unfollow"))
		return
	}

	hlog.Infof("Unfollow: userID=%d targetID=%d success", userID, req.TargetID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 获取粉丝列表
// @Description 获取用户的粉丝列表（需要认证）
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} FollowListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/followers [get]
func (h *ChatHandler) ListFollowers(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.chatClient.ListFollowers(ctx, userID)
	if err != nil {
		hlog.Errorf("ListFollowers: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get followers"))
		return
	}

	followers := make([]*UserItem, 0, len(resp.Followers))
	for _, f := range resp.Followers {
		followers = append(followers, &UserItem{
			ID:       f.Id,
			Username: f.Username,
			Avatar:   f.Avatar,
		})
	}

	c.JSON(200, hertz.Success(FollowListResponse{Followers: followers}))
}

// @Summary 获取关注列表
// @Description 获取用户关注的用户列表（需要认证）
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} FollowListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/following [get]
func (h *ChatHandler) ListFollowing(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.chatClient.ListFollowing(ctx, userID)
	if err != nil {
		hlog.Errorf("ListFollowing: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get following"))
		return
	}

	following := make([]*UserItem, 0, len(resp.Following))
	for _, f := range resp.Following {
		following = append(following, &UserItem{
			ID:       f.Id,
			Username: f.Username,
			Avatar:   f.Avatar,
		})
	}

	c.JSON(200, hertz.Success(FollowListResponse{Following: following}))
}

// @Summary 获取好友列表
// @Description 获取双向关注的好友列表（需要认证）
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} FriendListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/friends [get]
func (h *ChatHandler) ListFriends(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.chatClient.ListFriends(ctx, userID)
	if err != nil {
		hlog.Errorf("ListFriends: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get friends"))
		return
	}

	friends := make([]*FriendItem, 0, len(resp.Friends))
	for _, f := range resp.Friends {
		friends = append(friends, &FriendItem{
			ID:       f.Id,
			Username: f.Username,
			Avatar:   f.Avatar,
		})
	}

	c.JSON(200, hertz.Success(FriendListResponse{Friends: friends}))
}

// ========== 1v1 聊天 ==========

// @Summary 获取 WebSocket 连接（发起聊天）
// @Description 获取 WebSocket 连接地址和临时 token，用于进入聊天界面
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConnectRequest true "连接请求"
// @Success 200 {object} ConnectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "不是好友"
// @Failure 500 {object} ErrorResponse
// @Router /chat/connect [post]
func (h *ChatHandler) Connect(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req ConnectRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	peerID := req.PeerID

	// 验证好友关系（双向关注）
	isFriend, err := h.chatClient.IsFriends(ctx, userID, peerID)
	if err != nil {
		hlog.Errorf("Connect: check friend failed, userID=%d peerID=%d: %v", userID, peerID, err)
		c.JSON(500, hertz.Fail(500, "Failed to verify friend status"))
		return
	}
	if !isFriend {
		c.JSON(403, hertz.Fail(403, "Not friends, cannot chat"))
		return
	}

	// 生成 WebSocket token
	token, expiresAt, err := ws.GenerateWSToken(userID, peerID, "1v1", "chat-ws-secret")
	if err != nil {
		hlog.Errorf("Connect: generate token failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to generate token"))
		return
	}

	hlog.Infof("Connect: userID=%d peerID=%d success", userID, peerID)
	c.JSON(200, hertz.Success(ConnectResponse{
		WebSocketURL: fmt.Sprintf("ws://127.0.0.1:8892/ws/chat?token=%s", token),
		Token:        token,
		ExpiresAt:    expiresAt,
	}))
}

// @Summary 获取群聊 WebSocket 连接
// @Description 获取群聊的 WebSocket 连接地址和临时 token
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConnectGroupRequest true "群聊连接请求"
// @Success 200 {object} ConnectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "不是群成员"
// @Failure 500 {object} ErrorResponse
// @Router /chat/connect/group [post]
func (h *ChatHandler) ConnectGroup(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req ConnectGroupRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	groupID := req.GroupID

	// 验证群成员关系
	isMember, err := h.chatClient.IsGroupMember(ctx, userID, groupID)
	if err != nil {
		hlog.Errorf("ConnectGroup: check member failed, userID=%d groupID=%d: %v", userID, groupID, err)
		c.JSON(500, hertz.Fail(500, "Failed to verify group membership"))
		return
	}
	if !isMember {
		c.JSON(403, hertz.Fail(403, "Not a group member, cannot join"))
		return
	}

	// 生成 WebSocket token
	token, expiresAt, err := ws.GenerateWSToken(userID, groupID, "group", "chat-ws-secret")
	if err != nil {
		hlog.Errorf("ConnectGroup: generate token failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to generate token"))
		return
	}

	hlog.Infof("ConnectGroup: userID=%d groupID=%d success", userID, groupID)
	c.JSON(200, hertz.Success(ConnectResponse{
		WebSocketURL: fmt.Sprintf("ws://127.0.0.1:8892/ws/chat?token=%s", token),
		Token:        token,
		ExpiresAt:    expiresAt,
	}))
}

// @Summary 获取聊天记录
// @Description 获取与某个用户的聊天记录（需要认证）
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Param peer_id query int64 true "对方用户ID"
// @Param cursor query int64 false "cursor时间戳(毫秒)" default(0)
// @Param limit query int32 false "每页数量" default(10)
// @Success 200 {object} MessageListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/message/list [get]
func (h *ChatHandler) ListMessages(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	peerID, err := strconv.ParseInt(c.Query("peer_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid peer_id"))
		return
	}

	cursor, _ := strconv.ParseInt(c.DefaultQuery("cursor", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)

	resp, err := h.chatClient.ListMessages(ctx, userID, peerID, cursor, int32(limit))
	if err != nil {
		hlog.Errorf("ListMessages: userID=%d peerID=%d failed: %v", userID, peerID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get messages"))
		return
	}

	messages := make([]*MessageItem, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		messages = append(messages, &MessageItem{
			ID:             m.Id,
			ConversationID: m.ConversationId,
			SenderID:       m.SenderId,
			Content:        m.Content,
			CreatedAt:      m.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(MessageListResponse{
		Messages:   messages,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}))
}

// ========== 群聊 ==========

// @Summary 创建群聊
// @Description 创建群聊，创建者自动成为管理员（需要认证）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateGroupRequest true "创建群聊请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group [post]
func (h *ChatHandler) CreateGroup(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req CreateGroupRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.chatClient.CreateGroup(ctx, userID, req.Name)
	if err != nil {
		hlog.Errorf("CreateGroup: creatorID=%d name=%s failed: %v", userID, req.Name, err)
		c.JSON(500, hertz.Fail(500, "Failed to create group"))
		return
	}

	hlog.Infof("CreateGroup: creatorID=%d name=%s success, groupID=%d", userID, req.Name, resp.GroupId)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"group_id":        resp.GroupId,
		"conversation_id": resp.ConversationId,
		"success":         true,
	}))
}

// @Summary 添加群成员
// @Description 管理员添加群成员（需要认证，仅限管理员）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddGroupMemberRequest true "添加成员请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group/member [post]
func (h *ChatHandler) AddGroupMember(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req AddGroupMemberRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.chatClient.AddGroupMember(ctx, req.GroupID, req.UserID, userID)
	if err != nil {
		hlog.Errorf("AddGroupMember: adminID=%d groupID=%d userID=%d failed: %v", userID, req.GroupID, req.UserID, err)
		c.JSON(500, hertz.Fail(500, "Failed to add member"))
		return
	}

	hlog.Infof("AddGroupMember: adminID=%d groupID=%d userID=%d success", userID, req.GroupID, req.UserID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 移除群成员
// @Description 管理员移除群成员（需要认证，仅限管理员）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RemoveGroupMemberRequest true "移除成员请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group/member [delete]
func (h *ChatHandler) RemoveGroupMember(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req RemoveGroupMemberRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.chatClient.RemoveGroupMember(ctx, req.GroupID, req.UserID, userID)
	if err != nil {
		hlog.Errorf("RemoveGroupMember: adminID=%d groupID=%d userID=%d failed: %v", userID, req.GroupID, req.UserID, err)
		c.JSON(500, hertz.Fail(500, "Failed to remove member"))
		return
	}

	hlog.Infof("RemoveGroupMember: adminID=%d groupID=%d userID=%d success", userID, req.GroupID, req.UserID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 获取群成员列表
// @Description 获取群的所有成员列表
// @Tags chat
// @Produce json
// @Param group_id query int64 true "群ID"
// @Success 200 {object} GroupMemberListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group/members [get]
func (h *ChatHandler) ListGroupMembers(ctx context.Context, c *app.RequestContext) {
	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid group_id"))
		return
	}

	resp, err := h.chatClient.ListGroupMembers(ctx, groupID)
	if err != nil {
		hlog.Errorf("ListGroupMembers: groupID=%d failed: %v", groupID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get members"))
		return
	}

	members := make([]*GroupMemberItem, 0, len(resp.Members))
	for _, m := range resp.Members {
		members = append(members, &GroupMemberItem{
			UserID:   m.UserId,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	c.JSON(200, hertz.Success(GroupMemberListResponse{Members: members}))
}

// @Summary 获取我加入的群列表
// @Description 获取当前用户加入的所有群列表（需要认证）
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GroupListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/groups [get]
func (h *ChatHandler) ListUserGroups(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.chatClient.ListUserGroups(ctx, userID)
	if err != nil {
		hlog.Errorf("ListUserGroups: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get groups"))
		return
	}

	groups := make([]*GroupItem, 0, len(resp.Groups))
	for _, g := range resp.Groups {
		groups = append(groups, &GroupItem{
			ID:             g.Id,
			ConversationID: g.ConversationId,
			Name:           g.Name,
			Avatar:         g.Avatar,
			MemberCount:    g.MemberCount,
		})
	}

	c.JSON(200, hertz.Success(GroupListResponse{Groups: groups}))
}

// @Summary 发送群消息
// @Description 在群内发送消息（需要认证，仅限群成员）
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendGroupMessageRequest true "发送群消息请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group/message [post]
func (h *ChatHandler) SendGroupMessage(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req SendGroupMessageRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.chatClient.SendGroupMessage(ctx, userID, req.GroupID, req.Content)
	if err != nil {
		hlog.Errorf("SendGroupMessage: senderID=%d groupID=%d failed: %v", userID, req.GroupID, err)
		c.JSON(500, hertz.Fail(500, "Failed to send group message"))
		return
	}

	hlog.Infof("SendGroupMessage: senderID=%d groupID=%d success, messageID=%d", userID, req.GroupID, resp.MessageId)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"message_id":      resp.MessageId,
		"conversation_id": resp.ConversationId,
		"success":         true,
	}))
}

// @Summary 获取群消息历史
// @Description 获取群的消息历史记录
// @Tags chat
// @Produce json
// @Param group_id query int64 true "群ID"
// @Param cursor query int64 false "cursor时间戳(毫秒)" default(0)
// @Param limit query int32 false "每页数量" default(10)
// @Success 200 {object} MessageListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /chat/group/message/list [get]
func (h *ChatHandler) ListGroupMessages(ctx context.Context, c *app.RequestContext) {
	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid group_id"))
		return
	}

	cursor, _ := strconv.ParseInt(c.DefaultQuery("cursor", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)

	resp, err := h.chatClient.ListGroupMessages(ctx, groupID, cursor, int32(limit))
	if err != nil {
		hlog.Errorf("ListGroupMessages: groupID=%d failed: %v", groupID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get messages"))
		return
	}

	messages := make([]*MessageItem, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		messages = append(messages, &MessageItem{
			ID:             m.Id,
			ConversationID: m.ConversationId,
			SenderID:       m.SenderId,
			Content:        m.Content,
			CreatedAt:      m.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(MessageListResponse{
		Messages:   messages,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}))
}
