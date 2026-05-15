package router

import (
	"vicomova/internal/chat/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册聊天相关路由
// @Summary 聊天服务 HTTP API
// @Description 提供关注/粉丝、好友、1v1聊天、群聊的查询接口（消息发送通过 WebSocket）
// @Tags chat
// @Accept json
// @Produce json
func RegisterRoutes(h *server.Hertz, chatHandler *handler.ChatHandler) {
	g := h.Group("/chat", rootMw()...)

	// 关注相关
	follow := g.Group("/follow")
	// @Summary 关注用户
	// @Description 当前用户关注目标用户
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.FollowRequest true "关注请求"
	// @Success 200 {object} handler.SuccessResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/follow [post]
	follow.POST("", chatHandler.Follow)
	// @Summary 取消关注
	// @Description 当前用户取消关注目标用户
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.UnfollowRequest true "取消关注请求"
	// @Success 200 {object} handler.SuccessResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/follow [delete]
	follow.DELETE("", chatHandler.Unfollow)

	// @Summary 获取粉丝列表
	// @Description 获取指定用户的粉丝列表（关注该用户的用户）
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} handler.FollowListResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/followers [get]
	g.GET("/followers", chatHandler.ListFollowers)
	// @Summary 获取关注列表
	// @Description 获取指定用户关注的用户列表
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} handler.FollowListResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/following [get]
	g.GET("/following", chatHandler.ListFollowing)
	// @Summary 获取好友列表
	// @Description 获取指定用户的好友列表（双向关注即为好友）
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} handler.FriendListResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/friends [get]
	g.GET("/friends", chatHandler.ListFriends)

	// 1v1 聊天
	// @Summary 获取 WebSocket 连接（1v1聊天）
	// @Description 获取 WebSocket 连接地址和临时 token，用于进入聊天界面（需先验证好友关系）
	// @Tags chat
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.ConnectRequest true "连接请求（peer_id）"
	// @Success 200 {object} handler.ConnectResponse "返回 WebSocket URL 和临时 token"
	// @Failure 400 {object} handler.ErrorResponse "无效请求"
	// @Failure 401 {object} handler.ErrorResponse "未认证"
	// @Failure 403 {object} handler.ErrorResponse "不是好友，无法聊天"
	// @Failure 500 {object} handler.ErrorResponse "服务器错误"
	// @Router /chat/connect [post]
	g.POST("/connect", chatHandler.Connect)
	// @Summary 获取群聊 WebSocket 连接
	// @Description 获取群聊的 WebSocket 连接地址和临时 token（需验证群成员身份）
	// @Tags chat
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.ConnectGroupRequest true "群聊连接请求（group_id）"
	// @Success 200 {object} handler.ConnectResponse "返回 WebSocket URL 和临时 token"
	// @Failure 400 {object} handler.ErrorResponse "无效请求"
	// @Failure 401 {object} handler.ErrorResponse "未认证"
	// @Failure 403 {object} handler.ErrorResponse "不是群成员，无法加入"
	// @Failure 500 {object} handler.ErrorResponse "服务器错误"
	// @Router /chat/connect/group [post]
	g.POST("/connect/group", chatHandler.ConnectGroup)
	// @Summary 获取聊天记录
	// @Description 获取当前用户与指定用户的1v1聊天记录
	// @Produce json
	// @Security BearerAuth
	// @Param peer_id query int64 true "对方用户ID"
	// @Param cursor query int64 false "游标（时间戳毫秒）"
	// @Param limit query int false "每页数量，默认20"
	// @Success 200 {object} handler.MessageListResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/message/list [get]
	g.GET("/message/list", chatHandler.ListMessages)

	// 群聊
	// @Summary 创建群聊
	// @Description 创建一个新的群聊，当前用户为群主
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.CreateGroupRequest true "创建群聊请求"
	// @Success 200 {object} map[string]interface{} "返回 group_id, conversation_id"
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/group [post]
	g.POST("/group", chatHandler.CreateGroup)
	// @Summary 添加群成员
	// @Description 群主或管理员添加成员到群聊
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.AddGroupMemberRequest true "添加成员请求"
	// @Success 200 {object} handler.SuccessResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 403 {object} handler.ErrorResponse "不是管理员"
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/group/member [post]
	g.POST("/group/member", chatHandler.AddGroupMember)
	// @Summary 移除群成员
	// @Description 群主或管理员从群聊移除成员
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body handler.RemoveGroupMemberRequest true "移除成员请求"
	// @Success 200 {object} handler.SuccessResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 403 {object} handler.ErrorResponse "不是管理员"
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/group/member [delete]
	g.DELETE("/group/member", chatHandler.RemoveGroupMember)
	// @Summary 获取群成员列表
	// @Description 获取指定群聊的所有成员
	// @Produce json
	// @Param group_id query int64 true "群ID"
	// @Success 200 {object} handler.GroupMemberListResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/group/members [get]
	g.GET("/group/members", chatHandler.ListGroupMembers)
	// @Summary 获取我的群列表
	// @Description 获取当前用户加入的所有群聊
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} handler.GroupListResponse
	// @Failure 401 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/groups [get]
	g.GET("/groups", chatHandler.ListUserGroups)
	// @Summary 获取群消息记录
	// @Description 获取指定群聊的历史消息
	// @Produce json
	// @Param group_id query int64 true "群ID"
	// @Param cursor query int64 false "游标（时间戳毫秒）"
	// @Param limit query int false "每页数量，默认20"
	// @Success 200 {object} handler.MessageListResponse
	// @Failure 400 {object} handler.ErrorResponse
	// @Failure 500 {object} handler.ErrorResponse
	// @Router /chat/group/message/list [get]
	g.GET("/group/message/list", chatHandler.ListGroupMessages)
}
