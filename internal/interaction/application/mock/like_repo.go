package mock

import (
	"context"
	"fmt"

	"vicomova/internal/interaction/domain/entity"
)

type MockLikeRepository struct {
	Likes         map[string]*entity.Like // key: "userID-targetType-targetID"
	CreateErr     error
	DeleteErr     error
	GetErr        error
	ListByUserErr error
	CountErr      error
}

func NewMockLikeRepository() *MockLikeRepository {
	return &MockLikeRepository{
		Likes: make(map[string]*entity.Like),
	}
}

func (m *MockLikeRepository) key(userID int64, targetType string, targetID int64) string {
	return fmt.Sprintf("%d-%s-%d", userID, targetType, targetID)
}

func (m *MockLikeRepository) Create(ctx context.Context, like *entity.Like) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	like.ID = int64(len(m.Likes) + 1)
	m.Likes[m.key(like.UserID, like.TargetType, like.TargetID)] = like
	return nil
}

func (m *MockLikeRepository) Delete(ctx context.Context, userID int64, targetType string, targetID int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.Likes, m.key(userID, targetType, targetID))
	return nil
}

func (m *MockLikeRepository) Get(ctx context.Context, userID int64, targetType string, targetID int64) (*entity.Like, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Likes[m.key(userID, targetType, targetID)], nil
}

func (m *MockLikeRepository) ListByUser(ctx context.Context, userID int64, targetType string, cursor int64, limit int) ([]*entity.Like, bool, error) {
	if m.ListByUserErr != nil {
		return nil, false, m.ListByUserErr
	}
	var result []*entity.Like
	for _, l := range m.Likes {
		if l.UserID == userID && l.TargetType == targetType {
			result = append(result, l)
		}
	}
	if len(result) > limit {
		return result[:limit], true, nil
	}
	return result, false, nil
}

func (m *MockLikeRepository) Count(ctx context.Context, targetType string, targetID int64) (int64, error) {
	if m.CountErr != nil {
		return 0, m.CountErr
	}
	var count int64
	for _, l := range m.Likes {
		if l.TargetType == targetType && l.TargetID == targetID {
			count++
		}
	}
	return count, nil
}
