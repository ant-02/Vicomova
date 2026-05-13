package command

import (
	"context"
	"testing"

	"vicomova/internal/video/application/mock"
	"vicomova/internal/video/domain/entity"
	videoVO "vicomova/internal/video/domain/valueobject"
)

func newTestVideoCommandService(
	repo *mock.MockVideoRepository,
	cache *mock.MockVideoCache,
) *VideoCommandService {
	return NewVideoCommandService(repo, cache, nil)
}

// ===== Tests for Save =====

func TestSave_CreateNewVideo(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	svc := newTestVideoCommandService(repo, cache)

	cmd := &SaveVideoCommand{
		UserID:      1,
		Title:       "Test Video",
		Description: "Test Description",
		CategoryID:  1,
		CoverURL:    "http://example.com/cover.jpg",
		VideoURL:    "http://example.com/video.mp4",
		Duration:    120,
	}

	result, err := svc.Save(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.VideoID == 0 {
		t.Error("expected non-zero video ID")
	}

	// Verify video was created
	if len(repo.Videos) != 1 {
		t.Errorf("expected 1 video in repo, got %d", len(repo.Videos))
	}
}

func TestSave_UpdateExistingVideo(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	// Pre-create a video
	existingVideo := &entity.Video{
		ID:          100,
		UserID:      1,
		Title:       "Original Title",
		Description: "Original Description",
		CategoryID:  1,
		Status:      videoVO.VideoStatusEditing,
	}
	repo.Videos[100] = existingVideo

	svc := newTestVideoCommandService(repo, cache)

	cmd := &SaveVideoCommand{
		VideoID:     100,
		UserID:      1,
		Title:       "Updated Title",
		Description: "Updated Description",
		CategoryID:  2,
		CoverURL:    "http://example.com/new-cover.jpg",
		VideoURL:    "http://example.com/new-video.mp4",
		Duration:    200,
	}

	result, err := svc.Save(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.VideoID != 100 {
		t.Errorf("expected video ID 100, got %d", result.VideoID)
	}

	// Verify video was updated
	updated := repo.Videos[100]
	if updated.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", updated.Title)
	}
	if updated.CategoryID != 2 {
		t.Errorf("expected category ID 2, got %d", updated.CategoryID)
	}
	// Status should be reset to Editing
	if updated.Status != videoVO.VideoStatusEditing {
		t.Errorf("expected status Editing, got %v", updated.Status)
	}
	// Cache should be invalidated
	if _, ok := cache.Videos[100]; ok {
		t.Error("cache was not invalidated after update")
	}
}

func TestSave_VideoNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	svc := newTestVideoCommandService(repo, cache)

	cmd := &SaveVideoCommand{
		VideoID: 999,
		UserID:  1,
		Title:   "Test",
	}

	_, err := svc.Save(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when video not found")
	}
}

func TestSave_ForbiddenUpdateOtherUserVideo(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	existingVideo := &entity.Video{
		ID:     100,
		UserID: 1, // Owner is user 1
		Title:  "Original",
		Status: videoVO.VideoStatusEditing,
	}
	repo.Videos[100] = existingVideo

	svc := newTestVideoCommandService(repo, cache)

	cmd := &SaveVideoCommand{
		VideoID: 100,
		UserID:  2, // Different user
		Title:   "Hacked Title",
	}

	_, err := svc.Save(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when updating other user's video")
	}
}

// ===== Tests for Submit =====

func TestSubmit_Success(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	video := &entity.Video{
		ID:     100,
		UserID: 1,
		Title:  "Test Video",
		Status: videoVO.VideoStatusEditing,
	}
	repo.Videos[100] = video

	svc := newTestVideoCommandService(repo, cache)

	cmd := &SubmitVideoCommand{
		VideoID: 100,
		UserID:  1,
	}

	result, err := svc.Submit(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.VideoID != 100 {
		t.Errorf("expected video ID 100, got %d", result.VideoID)
	}
	if video.Status != videoVO.VideoStatusPending {
		t.Errorf("expected status Pending, got %v", video.Status)
	}
}

func TestSubmit_VideoNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	svc := newTestVideoCommandService(repo, cache)

	cmd := &SubmitVideoCommand{
		VideoID: 999,
		UserID:  1,
	}

	_, err := svc.Submit(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when video not found")
	}
}

func TestSubmit_ForbiddenSubmitOtherUserVideo(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	video := &entity.Video{
		ID:     100,
		UserID: 1,
		Status: videoVO.VideoStatusEditing,
	}
	repo.Videos[100] = video

	svc := newTestVideoCommandService(repo, cache)

	cmd := &SubmitVideoCommand{
		VideoID: 100,
		UserID:  2, // Different user
	}

	_, err := svc.Submit(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when submitting other user's video")
	}
}

func TestSubmit_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	video := &entity.Video{
		ID:     100,
		UserID: 1,
		Status: videoVO.VideoStatusPublished, // Already published, can't submit
	}
	repo.Videos[100] = video

	svc := newTestVideoCommandService(repo, cache)

	cmd := &SubmitVideoCommand{
		VideoID: 100,
		UserID:  1,
	}

	_, err := svc.Submit(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when video status is not Editing")
	}
}

// ===== Tests for Publish =====

func TestPublish_Success(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	video := &entity.Video{
		ID:     100,
		UserID: 1,
		Status: videoVO.VideoStatusPending,
	}
	repo.Videos[100] = video

	svc := newTestVideoCommandService(repo, cache)

	cmd := &PublishVideoCommand{
		VideoID: 100,
	}

	result, err := svc.Publish(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.VideoID != 100 {
		t.Errorf("expected video ID 100, got %d", result.VideoID)
	}
	if video.Status != videoVO.VideoStatusPublished {
		t.Errorf("expected status Published, got %v", video.Status)
	}
	// Cache should be invalidated
	if _, ok := cache.Videos[100]; ok {
		t.Error("cache was not invalidated after publish")
	}
}

func TestPublish_VideoNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()
	svc := newTestVideoCommandService(repo, cache)

	cmd := &PublishVideoCommand{
		VideoID: 999,
	}

	_, err := svc.Publish(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when video not found")
	}
}

func TestPublish_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	repo := mock.NewMockVideoRepository()
	cache := mock.NewMockVideoCache()

	video := &entity.Video{
		ID:     100,
		Status: videoVO.VideoStatusEditing, // Not pending, can't publish
	}
	repo.Videos[100] = video

	svc := newTestVideoCommandService(repo, cache)

	cmd := &PublishVideoCommand{
		VideoID: 100,
	}

	_, err := svc.Publish(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when video status is not Pending")
	}
}