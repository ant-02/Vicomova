package query

import (
	"context"
	"testing"
	"time"

	"vicomova/internal/interaction/application/mock"
	"vicomova/internal/interaction/domain/entity"
)

type mockErr string

func (e mockErr) Error() string { return string(e) }

var errInternal = mockErr("internal error")

func newTestService(
	likeRepo *mock.MockLikeRepository,
	commentRepo *mock.MockCommentRepository,
	favoriteRepo *mock.MockFavoriteRepository,
) *InteractionQueryService {
	return NewInteractionQueryService(likeRepo, commentRepo, favoriteRepo)
}

func TestListLikes_Success(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.Likes["1-video-100"] = &entity.Like{
		ID:         1,
		UserID:     1,
		TargetType: entity.TargetTypeVideo,
		TargetID:   100,
		CreatedAt:  time.Now(),
	}
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListLikes(context.Background(), 1, "video", 0, 10)
	if err != nil {
		t.Fatalf("ListLikes: expected nil, got %v", err)
	}
	if result == nil {
		t.Fatal("ListLikes: expected result, got nil")
	}
	_ = hasMore
}

func TestListLikes_Empty(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListLikes(context.Background(), 1, "video", 0, 10)
	if err != nil {
		t.Fatalf("ListLikes: expected nil, got %v", err)
	}
	if result.Items == nil {
		t.Fatal("ListLikes: expected items, got nil")
	}
	if hasMore {
		t.Error("ListLikes: expected hasMore=false")
	}
}

func TestListLikes_RepositoryError(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.ListByUserErr = errInternal
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	_, _, err := svc.ListLikes(context.Background(), 1, "video", 0, 10)
	if err == nil {
		t.Fatal("ListLikes: expected error, got nil")
	}
}

func TestListLikes_HasMore(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	for i := 0; i < 3; i++ {
		likeRepo.Likes["1-video-"+string(rune('0'+i))] = &entity.Like{
			ID:         int64(i),
			UserID:     1,
			TargetType: entity.TargetTypeVideo,
			TargetID:   int64(i),
			CreatedAt:  time.Now(),
		}
	}
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	_, hasMore, err := svc.ListLikes(context.Background(), 1, "video", 0, 2)
	if err != nil {
		t.Fatalf("ListLikes: expected nil, got %v", err)
	}
	if !hasMore {
		t.Error("ListLikes: expected hasMore=true")
	}
}

func TestListFavorites_Success(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.Favorites["1-100"] = &entity.Favorite{
		ID:        1,
		UserID:    1,
		VideoID:   100,
		CreatedAt: time.Now(),
	}
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	result, hasMore, err := svc.ListFavorites(context.Background(), 1, 0, 10)
	if err != nil {
		t.Fatalf("ListFavorites: expected nil, got %v", err)
	}
	if result == nil {
		t.Fatal("ListFavorites: expected result, got nil")
	}
	_ = hasMore
}

func TestListFavorites_Empty(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListFavorites(context.Background(), 1, 0, 10)
	if err != nil {
		t.Fatalf("ListFavorites: expected nil, got %v", err)
	}
	if result.Items == nil {
		t.Fatal("ListFavorites: expected items, got nil")
	}
	if hasMore {
		t.Error("ListFavorites: expected hasMore=false")
	}
}

func TestListFavorites_RepositoryError(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.ListByUserErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	_, _, err := svc.ListFavorites(context.Background(), 1, 0, 10)
	if err == nil {
		t.Fatal("ListFavorites: expected error, got nil")
	}
}

func TestListComments_RootComments(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[1] = &entity.Comment{
		ID:        1,
		UserID:    1,
		VideoID:   100,
		ParentID:  0,
		Content:   "root",
		CreatedAt: time.Now(),
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListComments(context.Background(), 100, 0, 0, 10)
	if err != nil {
		t.Fatalf("ListComments: expected nil, got %v", err)
	}
	if result == nil {
		t.Fatal("ListComments: expected result, got nil")
	}
	_ = hasMore
}

func TestListComments_Replies(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[1] = &entity.Comment{
		ID:       1,
		ParentID: 0,
	}
	commentRepo.Comments[2] = &entity.Comment{
		ID:       2,
		ParentID: 1,
		Content:  "reply",
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListComments(context.Background(), 100, 1, 0, 10)
	if err != nil {
		t.Fatalf("ListComments replies: expected nil, got %v", err)
	}
	if result == nil {
		t.Fatal("ListComments replies: expected result, got nil")
	}
	_ = hasMore
}

func TestListComments_Empty(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	result, hasMore, err := svc.ListComments(context.Background(), 100, 0, 0, 10)
	if err != nil {
		t.Fatalf("ListComments: expected nil, got %v", err)
	}
	if result.Items == nil {
		t.Fatal("ListComments: expected items, got nil")
	}
	if hasMore {
		t.Error("ListComments: expected hasMore=false")
	}
}

func TestListComments_ListByVideoError(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.ListByVideoErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	_, _, err := svc.ListComments(context.Background(), 100, 0, 0, 10)
	if err == nil {
		t.Fatal("ListComments: expected error, got nil")
	}
}

func TestListComments_ListByParentError(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.ListByParentErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	_, _, err := svc.ListComments(context.Background(), 100, 1, 0, 10)
	if err == nil {
		t.Fatal("ListComments replies: expected error, got nil")
	}
}

func TestIsLiked_True(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.Likes["1-video-100"] = &entity.Like{
		ID:         1,
		UserID:     1,
		TargetType: entity.TargetTypeVideo,
		TargetID:   100,
	}
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	liked, err := svc.IsLiked(context.Background(), 1, "video", 100)
	if err != nil {
		t.Fatalf("IsLiked: expected nil, got %v", err)
	}
	if !liked {
		t.Error("IsLiked: expected true")
	}
}

func TestIsLiked_False(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	liked, err := svc.IsLiked(context.Background(), 1, "video", 100)
	if err != nil {
		t.Fatalf("IsLiked: expected nil, got %v", err)
	}
	if liked {
		t.Error("IsLiked: expected false")
	}
}

func TestIsLiked_RepositoryError(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.GetErr = errInternal
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	_, err := svc.IsLiked(context.Background(), 1, "video", 100)
	if err == nil {
		t.Fatal("IsLiked: expected error, got nil")
	}
}

func TestIsFavorited_True(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.Favorites["1-100"] = &entity.Favorite{
		ID:      1,
		UserID:  1,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	favorited, err := svc.IsFavorited(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("IsFavorited: expected nil, got %v", err)
	}
	if !favorited {
		t.Error("IsFavorited: expected true")
	}
}

func TestIsFavorited_False(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	favorited, err := svc.IsFavorited(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("IsFavorited: expected nil, got %v", err)
	}
	if favorited {
		t.Error("IsFavorited: expected false")
	}
}

func TestIsFavorited_RepositoryError(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.GetErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	_, err := svc.IsFavorited(context.Background(), 1, 100)
	if err == nil {
		t.Fatal("IsFavorited: expected error, got nil")
	}
}
