package repository

import (
	"context"

	"vicomova/internal/chat/domain/entity"
)

type FollowRepository interface {
	Create(ctx context.Context, follow *entity.Follow) error
	Delete(ctx context.Context, followerID, followingID int64) error
	Get(ctx context.Context, followerID, followingID int64) (*entity.Follow, error)
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
	ListFollowers(ctx context.Context, userID int64) ([]int64, error)
	ListFollowing(ctx context.Context, userID int64) ([]int64, error)
	ListFriends(ctx context.Context, userID int64) ([]int64, error)
}

type ConversationRepository interface {
	Create(ctx context.Context, conv *entity.Conversation) error
	Create1v1(ctx context.Context, userID1, userID2 int64) (*entity.Conversation, error)
	Get1v1(ctx context.Context, userID1, userID2 int64) (*entity.Conversation, error)
	GetByID(ctx context.Context, id int64) (*entity.Conversation, error)
}

type GroupMemberRepository interface {
	Create(ctx context.Context, member *entity.GroupMember) error
	Delete(ctx context.Context, groupID, userID int64) error
	GetMember(ctx context.Context, groupID, userID int64) (*entity.GroupMember, error)
	IsMember(ctx context.Context, groupID, userID int64) (bool, error)
	IsAdmin(ctx context.Context, groupID, userID int64) (bool, error)
	ListByGroup(ctx context.Context, groupID int64) ([]*entity.GroupMember, error)
	ListByUser(ctx context.Context, userID int64) ([]int64, error)
}

type MessageRepository interface {
	Create(ctx context.Context, msg *entity.Message) error
	GetByID(ctx context.Context, id int64) (*entity.Message, error)
	ListByConversation(ctx context.Context, conversationID int64, cursor int64, limit int) ([]*entity.Message, bool, error)
}

type OfflineMessageRepository interface {
	Create(ctx context.Context, msg *entity.OfflineMessage) error
	ListByUser(ctx context.Context, userID int64, limit int) ([]*entity.OfflineMessage, error)
	MarkAsRead(ctx context.Context, userID int64) error
	UnreadCount(ctx context.Context, userID int64) (int64, error)
}
