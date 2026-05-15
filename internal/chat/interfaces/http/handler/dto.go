package handler

// ErrorResponse 错误响应
// @Summary 错误响应
// @Description API 错误时返回的响应格式
// @Tags common
type ErrorResponse struct {
	// @Description 错误码
	Code int32 `json:"code"`
	// @Description 错误信息
	Msg string `json:"msg"`
}

// FollowRequest 关注请求
// @Summary 关注请求
// @Description 用户关注另一个用户
// @Tags chat
type FollowRequest struct {
	// @Description 目标用户ID
	TargetID int64 `json:"target_id"`
}

// UnfollowRequest 取消关注请求
// @Summary 取消关注请求
// @Description 用户取消关注另一个用户
// @Tags chat
type UnfollowRequest struct {
	// @Description 目标用户ID
	TargetID int64 `json:"target_id"`
}

// SendMessageRequest 发送消息请求（已废弃，请使用 WebSocket）
// @Summary 发送消息请求
// @Description 向好友发送消息（建议使用 WebSocket）
// @Tags chat
type SendMessageRequest struct {
	// @Description 接收者用户ID
	ReceiverID int64 `json:"receiver_id"`
	// @Description 消息内容
	Content string `json:"content"`
}

// CreateGroupRequest 创建群聊请求
// @Summary 创建群聊请求
// @Description 创建一个新的群聊
// @Tags chat
type CreateGroupRequest struct {
	// @Description 群名称
	Name string `json:"name"`
}

// AddGroupMemberRequest 添加群成员请求
// @Summary 添加群成员请求
// @Description 管理员添加成员到群聊
// @Tags chat
type AddGroupMemberRequest struct {
	// @Description 群ID
	GroupID int64 `json:"group_id"`
	// @Description 被添加的用户ID
	UserID int64 `json:"user_id"`
}

// RemoveGroupMemberRequest 移除群成员请求
// @Summary 移除群成员请求
// @Description 管理员从群聊移除成员
// @Tags chat
type RemoveGroupMemberRequest struct {
	// @Description 群ID
	GroupID int64 `json:"group_id"`
	// @Description 被移除的用户ID
	UserID int64 `json:"user_id"`
}

// SendGroupMessageRequest 发送群消息请求（已废弃）
// @Summary 发送群消息请求
// @Description 在群内发送消息（建议使用 WebSocket）
// @Tags chat
type SendGroupMessageRequest struct {
	// @Description 群ID
	GroupID int64 `json:"group_id"`
	// @Description 消息内容
	Content string `json:"content"`
}

// UserItem 用户信息
// @Summary 用户信息
// @Description 用户基本信息
// @Tags chat
type UserItem struct {
	// @Description 用户ID
	ID int64 `json:"id"`
	// @Description 用户名
	Username string `json:"username"`
	// @Description 头像URL
	Avatar string `json:"avatar"`
}

// FollowListResponse 粉丝/关注列表响应
// @Summary 粉丝/关注列表响应
// @Description 获取粉丝或关注列表
// @Tags chat
type FollowListResponse struct {
	// @Description 粉丝列表
	Followers []*UserItem `json:"followers,omitempty"`
	// @Description 关注列表
	Following []*UserItem `json:"following,omitempty"`
}

// FriendItem 好友信息
// @Summary 好友信息
// @Description 好友基本信息
// @Tags chat
type FriendItem struct {
	// @Description 用户ID
	ID int64 `json:"id"`
	// @Description 用户名
	Username string `json:"username"`
	// @Description 头像URL
	Avatar string `json:"avatar"`
}

// FriendListResponse 好友列表响应
// @Summary 好友列表响应
// @Description 获取好友列表
// @Tags chat
type FriendListResponse struct {
	// @Description 好友列表
	Friends []*FriendItem `json:"friends"`
}

// MessageItem 消息信息
// @Summary 消息信息
// @Description 聊天消息详情
// @Tags chat
type MessageItem struct {
	// @Description 消息ID
	ID int64 `json:"id"`
	// @Description 会话ID
	ConversationID int64 `json:"conversation_id"`
	// @Description 发送者用户ID
	SenderID int64 `json:"sender_id"`
	// @Description 消息内容
	Content string `json:"content"`
	// @Description 创建时间戳（秒）
	CreatedAt int64 `json:"created_at"`
}

// MessageListResponse 消息列表响应
// @Summary 消息列表响应
// @Description 获取聊天记录列表
// @Tags chat
type MessageListResponse struct {
	// @Description 消息列表
	Messages []*MessageItem `json:"messages"`
	// @Description 下页游标
	NextCursor string `json:"next_cursor"`
	// @Description 是否有更多
	HasMore bool `json:"has_more"`
}

// GroupMemberItem 群成员信息
// @Summary 群成员信息
// @Description 群成员详情
// @Tags chat
type GroupMemberItem struct {
	// @Description 用户ID
	UserID int64 `json:"user_id"`
	// @Description 角色（0=成员, 1=管理员）
	Role int64 `json:"role"`
	// @Description 加入时间戳（秒）
	JoinedAt int64 `json:"joined_at"`
}

// GroupMemberListResponse 群成员列表响应
// @Summary 群成员列表响应
// @Description 获取群成员列表
// @Tags chat
type GroupMemberListResponse struct {
	// @Description 群成员列表
	Members []*GroupMemberItem `json:"members"`
}

// GroupItem 群聊信息
// @Summary 群聊信息
// @Description 群聊基本信息
// @Tags chat
type GroupItem struct {
	// @Description 群ID
	ID int64 `json:"id"`
	// @Description 会话ID
	ConversationID int64 `json:"conversation_id"`
	// @Description 群名称
	Name string `json:"name"`
	// @Description 群头像
	Avatar string `json:"avatar"`
	// @Description 成员数量
	MemberCount int64 `json:"member_count"`
}

// GroupListResponse 群列表响应
// @Summary 群列表响应
// @Description 获取用户加入的群列表
// @Tags chat
type GroupListResponse struct {
	// @Description 群列表
	Groups []*GroupItem `json:"groups"`
}

// SuccessResponse 成功响应
// @Summary 成功响应
// @Description 操作成功响应
// @Tags common
type SuccessResponse struct {
	// @Description 是否成功
	Success bool `json:"success"`
	// @Description 消息（可选）
	Message string `json:"message,omitempty"`
}

// ConnectRequest 发起 WebSocket 连接请求（1v1聊天）
// @Summary 发起 WebSocket 连接请求（1v1）
// @Description 获取 1v1 聊天的 WebSocket 连接
// @Tags chat
type ConnectRequest struct {
	// @Description 对方用户ID
	PeerID int64 `json:"peer_id"`
}

// ConnectGroupRequest 发起 WebSocket 连接请求（群聊）
// @Summary 发起 WebSocket 连接请求（群聊）
// @Description 获取群聊的 WebSocket 连接
// @Tags chat
type ConnectGroupRequest struct {
	// @Description 群ID
	GroupID int64 `json:"group_id"`
}

// ConnectResponse WebSocket 连接响应
// @Summary WebSocket 连接响应
// @Description 返回 WebSocket 连接地址和临时 token
// @Tags chat
type ConnectResponse struct {
	// @Description WebSocket 连接地址
	WebSocketURL string `json:"websocket_url"`
	// @Description 临时连接 token
	Token string `json:"token"`
	// @Description 过期时间戳（秒）
	ExpiresAt int64 `json:"expires_at"`
}
