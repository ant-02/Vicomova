package mock

import (
	"context"
	"fmt"

	"vicomova/internal/interaction/domain/entity"
)

type MockFavoriteRepository struct {
	Favorites     map[string]*entity.Favorite // key: "userID-videoID"
	CreateErr     error
	DeleteErr     error
	GetErr        error
	ListByUserErr error
}

func NewMockFavoriteRepository() *MockFavoriteRepository {
	return &MockFavoriteRepository{
		Favorites: make(map[string]*entity.Favorite),
	}
}

func (m *MockFavoriteRepository) key(userID, videoID int64) string {
	return fmt.Sprintf("%d-%d", userID, videoID)
}

func (m *MockFavoriteRepository) Create(ctx context.Context, fav *entity.Favorite) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	fav.ID = int64(len(m.Favorites) + 1)
	m.Favorites[m.key(fav.UserID, fav.VideoID)] = fav
	return nil
}

func (m *MockFavoriteRepository) Delete(ctx context.Context, userID, videoID int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.Favorites, m.key(userID, videoID))
	return nil
}

func (m *MockFavoriteRepository) Get(ctx context.Context, userID, videoID int64) (*entity.Favorite, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Favorites[m.key(userID, videoID)], nil
}

func (m *MockFavoriteRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*entity.Favorite, bool, error) {
	if m.ListByUserErr != nil {
		return nil, false, m.ListByUserErr
	}
	var result []*entity.Favorite
	for _, f := range m.Favorites {
		if f.UserID == userID {
			result = append(result, f)
		}
	}
	if len(result) > limit {
		return result[:limit], true, nil
	}
	return result, false, nil
}
