package grpc

import (
	"context"

	chat "vicomova/third_party/kitex_gen/chat"
	chatservice "vicomova/third_party/kitex_gen/chat/chatservice"

	"github.com/cloudwego/kitex/client"
)

type ChatClient struct {
	cli chatservice.Client
}

// NewChatClient creates RPC client, directly connects to specified address
func NewChatClient(serviceName, addr string) (*ChatClient, error) {
	cli, err := chatservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &ChatClient{cli: cli}, nil
}

// ========== 关注相关 ==========

func (c *ChatClient) Follow(ctx context.Context, userID, targetID int64) (*chat.FollowResponse, error) {
	return c.cli.Follow(ctx, &chat.FollowRequest{
		UserId:   userID,
		TargetId: targetID,
	})
}

func (c *ChatClient) Unfollow(ctx context.Context, userID, targetID int64) (*chat.UnfollowResponse, error) {
	return c.cli.Unfollow(ctx, &chat.UnfollowRequest{
		UserId:   userID,
		TargetId: targetID,
	})
}

func (c *ChatClient) ListFollowers(ctx context.Context, userID int64) (*chat.ListFollowersResponse, error) {
	return c.cli.ListFollowers(ctx, &chat.ListFollowersRequest{
		UserId: userID,
	})
}

func (c *ChatClient) ListFollowing(ctx context.Context, userID int64) (*chat.ListFollowingResponse, error) {
	return c.cli.ListFollowing(ctx, &chat.ListFollowingRequest{
		UserId: userID,
	})
}

func (c *ChatClient) ListFriends(ctx context.Context, userID int64) (*chat.ListFriendsResponse, error) {
	return c.cli.ListFriends(ctx, &chat.ListFriendsRequest{
		UserId: userID,
	})
}

func (c *ChatClient) IsFriends(ctx context.Context, userID, targetID int64) (bool, error) {
	resp, err := c.cli.IsFriends(ctx, &chat.IsFriendsRequest{
		UserId:   userID,
		TargetId: targetID,
	})
	if err != nil {
		return false, err
	}
	return resp.IsFriends, nil
}

func (c *ChatClient) IsGroupMember(ctx context.Context, userID, groupID int64) (bool, error) {
	resp, err := c.cli.IsGroupMember(ctx, &chat.IsGroupMemberRequest{
		UserId:  userID,
		GroupId: groupID,
	})
	if err != nil {
		return false, err
	}
	return resp.IsMember, nil
}

// ========== 1v1 聊天 ==========

func (c *ChatClient) SendMessage(ctx context.Context, senderID, receiverID int64, content string) (*chat.SendMessageResponse, error) {
	return c.cli.SendMessage(ctx, &chat.SendMessageRequest{
		SenderId:   senderID,
		ReceiverId: receiverID,
		Content:    content,
	})
}

func (c *ChatClient) ListMessages(ctx context.Context, userID, peerID int64, cursor int64, limit int32) (*chat.ListMessagesResponse, error) {
	return c.cli.ListMessages(ctx, &chat.ListMessagesRequest{
		UserId: userID,
		PeerId: peerID,
		Cursor: cursor,
		Limit:  limit,
	})
}

// ========== 群聊 ==========

func (c *ChatClient) CreateGroup(ctx context.Context, creatorID int64, name string) (*chat.CreateGroupResponse, error) {
	return c.cli.CreateGroup(ctx, &chat.CreateGroupRequest{
		CreatorId: creatorID,
		Name:      name,
	})
}

func (c *ChatClient) AddGroupMember(ctx context.Context, groupID, userID, adminID int64) (*chat.AddGroupMemberResponse, error) {
	return c.cli.AddGroupMember(ctx, &chat.AddGroupMemberRequest{
		GroupId: groupID,
		UserId:  userID,
		AdminId: adminID,
	})
}

func (c *ChatClient) RemoveGroupMember(ctx context.Context, groupID, userID, adminID int64) (*chat.RemoveGroupMemberResponse, error) {
	return c.cli.RemoveGroupMember(ctx, &chat.RemoveGroupMemberRequest{
		GroupId: groupID,
		UserId:  userID,
		AdminId: adminID,
	})
}

func (c *ChatClient) ListGroupMembers(ctx context.Context, groupID int64) (*chat.ListGroupMembersResponse, error) {
	return c.cli.ListGroupMembers(ctx, &chat.ListGroupMembersRequest{
		GroupId: groupID,
	})
}

func (c *ChatClient) ListUserGroups(ctx context.Context, userID int64) (*chat.ListUserGroupsResponse, error) {
	return c.cli.ListUserGroups(ctx, &chat.ListUserGroupsRequest{
		UserId: userID,
	})
}

func (c *ChatClient) SendGroupMessage(ctx context.Context, senderID, groupID int64, content string) (*chat.SendGroupMessageResponse, error) {
	return c.cli.SendGroupMessage(ctx, &chat.SendGroupMessageRequest{
		SenderId: senderID,
		GroupId:  groupID,
		Content:  content,
	})
}

func (c *ChatClient) ListGroupMessages(ctx context.Context, groupID int64, cursor int64, limit int32) (*chat.ListGroupMessagesResponse, error) {
	return c.cli.ListGroupMessages(ctx, &chat.ListGroupMessagesRequest{
		GroupId: groupID,
		Cursor:  cursor,
		Limit:   limit,
	})
}
