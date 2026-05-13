package mock

import (
	"context"

	"vicomova/internal/video/domain/entity"
	videoVO "vicomova/internal/video/domain/valueobject"
)

type MockVideoRepository struct {
	Videos       map[int64]*entity.Video
	CreateErr    error
	GetByIDErr   error
	UpdateErr    error
	DeleteErr    error
	ListHotErr   error
	ListByCatErr error
}

func NewMockVideoRepository() *MockVideoRepository {
	return &MockVideoRepository{
		Videos: make(map[int64]*entity.Video),
	}
}

func (m *MockVideoRepository) Create(ctx context.Context, v *entity.Video) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	v.ID = int64(len(m.Videos) + 1)
	m.Videos[v.ID] = v
	return nil
}

func (m *MockVideoRepository) GetByID(ctx context.Context, id int64) (*entity.Video, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.Videos[id], nil
}

func (m *MockVideoRepository) GetCounts(ctx context.Context, videoID int64) (likeCount, viewCount int64, err error) {
	v := m.Videos[videoID]
	if v == nil {
		return 0, 0, nil
	}
	return v.LikeCount, v.ViewCount, nil
}

func (m *MockVideoRepository) Update(ctx context.Context, v *entity.Video) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	m.Videos[v.ID] = v
	return nil
}

func (m *MockVideoRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.Videos, id)
	return nil
}

func (m *MockVideoRepository) ListByCategory(ctx context.Context, categoryID int, page, size int) ([]*entity.Video, int64, error) {
	if m.ListByCatErr != nil {
		return nil, 0, m.ListByCatErr
	}
	var result []*entity.Video
	for _, v := range m.Videos {
		if v.CategoryID == categoryID && v.Status == videoVO.VideoStatusPublished {
			result = append(result, v)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockVideoRepository) ListByUser(ctx context.Context, userID int64, page, size int) ([]*entity.Video, int64, error) {
	var allVideos []*entity.Video
	for _, v := range m.Videos {
		if v.UserID == userID {
			allVideos = append(allVideos, v)
		}
	}
	total := int64(len(allVideos))
	if len(allVideos) == 0 {
		return allVideos, total, nil
	}
	// Apply pagination
	start := (page - 1) * size
	if start >= len(allVideos) {
		return []*entity.Video{}, total, nil
	}
	end := start + size
	if end > len(allVideos) {
		end = len(allVideos)
	}
	return allVideos[start:end], total, nil
}

func (m *MockVideoRepository) ListPublished(ctx context.Context, page, size int) ([]*entity.Video, int64, error) {
	var result []*entity.Video
	for _, v := range m.Videos {
		if v.Status == videoVO.VideoStatusPublished {
			result = append(result, v)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockVideoRepository) ListHot(ctx context.Context, limit int) ([]*entity.Video, error) {
	if m.ListHotErr != nil {
		return nil, m.ListHotErr
	}
	result := make([]*entity.Video, 0, limit)
	for _, v := range m.Videos {
		if v.Status == videoVO.VideoStatusPublished {
			result = append(result, v)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockVideoRepository) IncrementView(ctx context.Context, id int64) error {
	if v := m.Videos[id]; v != nil {
		v.ViewCount++
	}
	return nil
}

func (m *MockVideoRepository) IncrementViewBatch(ctx context.Context, id int64, delta int64) error {
	if v := m.Videos[id]; v != nil {
		v.ViewCount += delta
	}
	return nil
}

func (m *MockVideoRepository) UpdateCounts(ctx context.Context, id int64, likeDelta, commentDelta int64) error {
	if v := m.Videos[id]; v != nil {
		v.LikeCount += likeDelta
		v.CommentCount += commentDelta
	}
	return nil
}
