package mock

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

type MockVideoCache struct {
	Videos map[int64]*entity.Video
	GetErr error
	SetErr error
	DelErr error
}

func NewMockVideoCache() *MockVideoCache {
	return &MockVideoCache{
		Videos: make(map[int64]*entity.Video),
	}
}

func (m *MockVideoCache) Get(ctx context.Context, id int64) (*entity.Video, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Videos[id], nil
}

func (m *MockVideoCache) GetAndRefresh(ctx context.Context, id int64) (*entity.Video, error) {
	return m.Get(ctx, id)
}

func (m *MockVideoCache) Set(ctx context.Context, video *entity.Video) error {
	if m.SetErr != nil {
		return m.SetErr
	}
	m.Videos[video.ID] = video
	return nil
}

func (m *MockVideoCache) Del(ctx context.Context, id int64) error {
	if m.DelErr != nil {
		return m.DelErr
	}
	delete(m.Videos, id)
	return nil
}

type MockHotVideoCache struct {
	Metas      map[int64]*entity.HotVideoMeta
	HotVideos  []entity.HotVideoScore
	GetMetaErr error
	SetMetaErr error
	GetIDsErr  error
	SetIDsErr  error
	CountErr   error
}

func NewMockHotVideoCache() *MockHotVideoCache {
	return &MockHotVideoCache{
		Metas: make(map[int64]*entity.HotVideoMeta),
	}
}

func (m *MockHotVideoCache) GetHotVideoIDs(ctx context.Context, offset, limit int64) ([]entity.HotVideoScore, error) {
	if m.GetIDsErr != nil {
		return nil, m.GetIDsErr
	}
	if int64(len(m.HotVideos)) <= offset {
		return nil, nil
	}
	end := offset + limit
	if end > int64(len(m.HotVideos)) {
		end = int64(len(m.HotVideos))
	}
	return m.HotVideos[offset:end], nil
}

func (m *MockHotVideoCache) GetHotVideoCount(ctx context.Context) (int64, error) {
	if m.CountErr != nil {
		return 0, m.CountErr
	}
	return int64(len(m.HotVideos)), nil
}

func (m *MockHotVideoCache) SetHotVideoIDs(ctx context.Context, videos []*entity.Video) error {
	return m.SetIDsErr
}

func (m *MockHotVideoCache) SetHotVideoIDsStaging(ctx context.Context, videos []*entity.Video) error {
	return nil
}

func (m *MockHotVideoCache) SwapHotVideoIDs(ctx context.Context) error {
	return nil
}

func (m *MockHotVideoCache) GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error) {
	if m.GetMetaErr != nil {
		return nil, m.GetMetaErr
	}
	return m.Metas[videoID], nil
}

func (m *MockHotVideoCache) SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error {
	if m.SetMetaErr != nil {
		return m.SetMetaErr
	}
	m.Metas[videoID] = meta
	return nil
}

func (m *MockHotVideoCache) SetHotVideo(ctx context.Context, videoID int64, score float64) error {
	return nil
}

func (m *MockHotVideoCache) RemoveHotVideo(ctx context.Context, videoID int64) error {
	return nil
}

func (m *MockHotVideoCache) ClearHotVideos(ctx context.Context) error {
	m.HotVideos = nil
	return nil
}

type MockCategoryVideoCache struct {
	Metas  map[int64]*entity.HotVideoMeta
	Videos map[int][]entity.HotVideoScore
	GetErr error
	SetErr error
}

func NewMockCategoryVideoCache() *MockCategoryVideoCache {
	return &MockCategoryVideoCache{
		Metas:  make(map[int64]*entity.HotVideoMeta),
		Videos: make(map[int][]entity.HotVideoScore),
	}
}

func (m *MockCategoryVideoCache) GetVideoIDs(ctx context.Context, categoryID int, offset, limit int64) ([]entity.HotVideoScore, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	vids, ok := m.Videos[categoryID]
	if !ok || int64(len(vids)) <= offset {
		return nil, nil
	}
	end := offset + limit
	if end > int64(len(vids)) {
		end = int64(len(vids))
	}
	return vids[offset:end], nil
}

func (m *MockCategoryVideoCache) GetVideoCount(ctx context.Context, categoryID int) (int64, error) {
	return int64(len(m.Videos[categoryID])), nil
}

func (m *MockCategoryVideoCache) GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Metas[videoID], nil
}

func (m *MockCategoryVideoCache) SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error {
	if m.SetErr != nil {
		return m.SetErr
	}
	m.Metas[videoID] = meta
	return nil
}

func (m *MockCategoryVideoCache) SetVideoIDs(ctx context.Context, categoryID int, videos []*entity.Video) error {
	if m.SetErr != nil {
		return m.SetErr
	}
	scores := make([]entity.HotVideoScore, len(videos))
	for i, v := range videos {
		scores[i] = entity.HotVideoScore{VideoID: v.ID, Score: v.HotScore}
	}
	m.Videos[categoryID] = scores
	return nil
}

func (m *MockCategoryVideoCache) ClearCategoryVideos(ctx context.Context, categoryID int) error {
	delete(m.Videos, categoryID)
	return nil
}

type MockViewCountProducer struct {
	RecordCalled bool
	RecordErr    error
	VideoIDs     []int64
}

func NewMockViewCountProducer() *MockViewCountProducer {
	return &MockViewCountProducer{}
}

func (m *MockViewCountProducer) Record(ctx context.Context, videoID int64) error {
	if m.RecordErr != nil {
		return m.RecordErr
	}
	m.RecordCalled = true
	m.VideoIDs = append(m.VideoIDs, videoID)
	return nil
}

func (m *MockViewCountProducer) RecordBatch(ctx context.Context, videoIDs []int64) error {
	if m.RecordErr != nil {
		return m.RecordErr
	}
	m.VideoIDs = append(m.VideoIDs, videoIDs...)
	return nil
}
