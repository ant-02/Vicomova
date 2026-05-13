package query

import (
	"context"
	"testing"

	"vicomova/internal/video/application/mock"
	"vicomova/internal/video/domain/entity"
	videoVO "vicomova/internal/video/domain/valueobject"
)

func newTestVideoQueryService(
	repo *mock.MockVideoRepository,
	cache *mock.MockVideoCache,
	hotCache *mock.MockHotVideoCache,
	categoryCache *mock.MockCategoryVideoCache,
	viewCountProducer *mock.MockViewCountProducer,
) *VideoQueryService {
	return NewVideoQueryService(
		repo,
		cache,
		nil, // oss
		viewCountProducer,
		hotCache,
		categoryCache,
		nil, // userClient
	)
}

// ===== Tests for GetVideoStream =====

func TestGetVideoStream_CacheHit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	video := &entity.Video{
		ID:        100,
		UserID:    1,
		Title:     "Cached Video",
		Status:    videoVO.VideoStatusPublished,
		ViewCount: 100,
	}
	cache.Videos[100] = video

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.GetVideoStream(ctx, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Video.Title != "Cached Video" {
		t.Errorf("expected title 'Cached Video', got '%s'", result.Video.Title)
	}
	// View count should be recorded
	if !producer.RecordCalled {
		t.Error("expected view count to be recorded")
	}
}

func TestGetVideoStream_CacheMiss_DBHit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	video := &entity.Video{
		ID:        100,
		UserID:    1,
		Title:     "DB Video",
		Status:    videoVO.VideoStatusPublished,
		ViewCount: 100,
	}
	repo.Videos[100] = video

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.GetVideoStream(ctx, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Video.Title != "DB Video" {
		t.Errorf("expected title 'DB Video', got '%s'", result.Video.Title)
	}
	// Video should be cached
	if _, ok := cache.Videos[100]; !ok {
		t.Error("expected video to be cached")
	}
}

func TestGetVideoStream_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	_, err := svc.GetVideoStream(ctx, 999)
	if err == nil {
		t.Fatal("expected error when video not found")
	}
}

func TestGetVideoStream_NotPublished(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	video := &entity.Video{
		ID:     100,
		UserID: 1,
		Status: videoVO.VideoStatusEditing, // Not published
	}
	cache.Videos[100] = video

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	_, err := svc.GetVideoStream(ctx, 100)
	if err == nil {
		t.Fatal("expected error when video is not published")
	}
}

// ===== Tests for ListPublishedVideos =====

func TestListPublishedVideos_Success(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	repo.Videos[1] = &entity.Video{ID: 1, UserID: 1, Title: "Video 1", Status: videoVO.VideoStatusPublished}
	repo.Videos[2] = &entity.Video{ID: 2, UserID: 1, Title: "Video 2", Status: videoVO.VideoStatusPublished}
	repo.Videos[3] = &entity.Video{ID: 3, UserID: 2, Title: "Video 3", Status: videoVO.VideoStatusPublished}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListPublishedVideos(ctx, 1, &PublishedVideosQuery{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(result.Videos))
	}
	if result.HasMore {
		t.Error("expected hasMore=false")
	}
}

func TestListPublishedVideos_WithCursor(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	for i := 1; i <= 5; i++ {
		repo.Videos[int64(i)] = &entity.Video{ID: int64(i), UserID: 1, Title: "Video", Status: videoVO.VideoStatusPublished}
	}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListPublishedVideos(ctx, 1, &PublishedVideosQuery{Limit: 2, Cursor: "2"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(result.Videos))
	}
}

// ===== Tests for ListHotVideos =====

func TestListHotVideos_CacheHit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	// Setup hot cache with video IDs
	hotCache.HotVideos = []entity.HotVideoScore{
		{VideoID: 100, Score: 100},
		{VideoID: 200, Score: 90},
	}
	// Setup video metas
	hotCache.Metas[100] = &entity.HotVideoMeta{
		VideoID: 100, Title: "Hot Video 1", ViewCount: 1000, CommentCount: 50,
	}
	hotCache.Metas[200] = &entity.HotVideoMeta{
		VideoID: 200, Title: "Hot Video 2", ViewCount: 900, CommentCount: 40,
	}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListHotVideos(ctx, &ListHotVideosQuery{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(result.Videos))
	}
}

func TestListHotVideos_CacheMiss_FallbackToRepo(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	// Add videos directly to repo (no cache)
	repo.Videos[1] = &entity.Video{
		ID: 1, Title: "Fallback Video 1", Status: videoVO.VideoStatusPublished,
		ViewCount: 100, CommentCount: 10,
	}
	repo.Videos[2] = &entity.Video{
		ID: 2, Title: "Fallback Video 2", Status: videoVO.VideoStatusPublished,
		ViewCount: 90, CommentCount: 8,
	}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListHotVideos(ctx, &ListHotVideosQuery{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should fallback to repo and return videos
	if len(result.Videos) == 0 {
		t.Error("expected non-empty videos from fallback")
	}
}

func TestListHotVideos_InvalidLimit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListHotVideos(ctx, &ListHotVideosQuery{Limit: -1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should use default limit
	if result.Videos == nil {
		t.Error("expected non-nil result")
	}
}

// ===== Tests for ListCategoryVideos =====

func TestListCategoryVideos_CacheHit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	// Setup category cache
	categoryCache.Videos[1] = []entity.HotVideoScore{
		{VideoID: 100, Score: 100},
		{VideoID: 200, Score: 90},
	}
	categoryCache.Metas[100] = &entity.HotVideoMeta{
		VideoID: 100, Title: "Category Video 1", ViewCount: 500, CommentCount: 20,
	}
	categoryCache.Metas[200] = &entity.HotVideoMeta{
		VideoID: 200, Title: "Category Video 2", ViewCount: 400, CommentCount: 15,
	}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListCategoryVideos(ctx, &CategoryVideosQuery{CategoryID: 1, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 2 {
		t.Errorf("expected 2 videos, got %d", len(result.Videos))
	}
}

func TestListCategoryVideos_CacheMiss_DBHit(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	// Setup videos in repo
	repo.Videos[1] = &entity.Video{
		ID:         1,
		UserID:     1,
		Title:      "DB Category Video",
		CategoryID: 1,
		Status:     videoVO.VideoStatusPublished,
		ViewCount:  100,
	}

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListCategoryVideos(ctx, &CategoryVideosQuery{CategoryID: 1, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 1 {
		t.Errorf("expected 1 video, got %d", len(result.Videos))
	}
}

func TestListCategoryVideos_EmptyResult(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	hotCache := mock.NewMockHotVideoCache()
	categoryCache := mock.NewMockCategoryVideoCache()
	producer := mock.NewMockViewCountProducer()

	svc := newTestVideoQueryService(repo, cache, hotCache, categoryCache, producer)

	result, err := svc.ListCategoryVideos(ctx, &CategoryVideosQuery{CategoryID: 999, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Videos) != 0 {
		t.Errorf("expected 0 videos, got %d", len(result.Videos))
	}
}
