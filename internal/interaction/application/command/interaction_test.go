package command

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
) *InteractionCommandService {
	return NewInteractionCommandService(likeRepo, commentRepo, favoriteRepo)
}

func TestLikeVideo_Success(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.LikeVideo(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err != nil {
		t.Fatalf("LikeVideo: expected nil, got %v", err)
	}
}

func TestLikeVideo_AlreadyLiked(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.Likes["1-video-100"] = &entity.Like{
		ID:         1,
		UserID:     1,
		TargetType: entity.TargetTypeVideo,
		TargetID:   100,
	}
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.LikeVideo(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err != nil {
		t.Fatalf("LikeVideo: expected nil, got %v", err)
	}
}

func TestLikeVideo_RepositoryError(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.GetErr = errInternal
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.LikeVideo(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err == nil {
		t.Fatal("LikeVideo: expected error, got nil")
	}
}

func TestUnlikeVideo_Success(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.Likes["1-video-100"] = &entity.Like{
		ID:         1,
		UserID:     1,
		TargetType: entity.TargetTypeVideo,
		TargetID:   100,
	}
	svc := newTestService(likeRepo, mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.UnlikeVideo(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err != nil {
		t.Fatalf("UnlikeVideo: expected nil, got %v", err)
	}
	if likeRepo.Likes["1-video-100"] != nil {
		t.Error("UnlikeVideo: like should be deleted")
	}
}

func TestLikeComment_Success(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[100] = &entity.Comment{
		ID:        100,
		UserID:    2,
		VideoID:   1,
		ParentID:  0,
		Content:   "test",
		LikeCount: 0,
		CreatedAt: time.Now(),
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	err := svc.LikeComment(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err != nil {
		t.Fatalf("LikeComment: expected nil, got %v", err)
	}
	if commentRepo.Comments[100].LikeCount != 1 {
		t.Errorf("LikeComment: expected like_count=1, got %d", commentRepo.Comments[100].LikeCount)
	}
}

func TestLikeComment_AlreadyLiked(t *testing.T) {
	likeRepo := mock.NewMockLikeRepository()
	likeRepo.Likes["1-comment-100"] = &entity.Like{
		ID:         1,
		UserID:     1,
		TargetType: entity.TargetTypeComment,
		TargetID:   100,
	}
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[100] = &entity.Comment{
		ID:        100,
		LikeCount: 0,
	}
	svc := newTestService(likeRepo, commentRepo, mock.NewMockFavoriteRepository())

	err := svc.LikeComment(context.Background(), &LikeCommand{
		UserID:   1,
		TargetID: 100,
	})
	if err != nil {
		t.Fatalf("LikeComment: expected nil, got %v", err)
	}
	if commentRepo.Comments[100].LikeCount != 0 {
		t.Error("LikeComment: should not increment when already liked")
	}
}

func TestAddFavorite_Success(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.AddFavorite(context.Background(), &FavoriteCommand{
		UserID:  1,
		VideoID: 100,
	})
	if err != nil {
		t.Fatalf("AddFavorite: expected nil, got %v", err)
	}
}

func TestAddFavorite_AlreadyFavorited(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.Favorites["1-100"] = &entity.Favorite{
		ID:      1,
		UserID:  1,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	err := svc.AddFavorite(context.Background(), &FavoriteCommand{
		UserID:  1,
		VideoID: 100,
	})
	if err != nil {
		t.Fatalf("AddFavorite: expected nil, got %v", err)
	}
}

func TestRemoveFavorite_Success(t *testing.T) {
	favoriteRepo := mock.NewMockFavoriteRepository()
	favoriteRepo.Favorites["1-100"] = &entity.Favorite{
		ID:      1,
		UserID:  1,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), favoriteRepo)

	err := svc.RemoveFavorite(context.Background(), &FavoriteCommand{
		UserID:  1,
		VideoID: 100,
	})
	if err != nil {
		t.Fatalf("RemoveFavorite: expected nil, got %v", err)
	}
	if favoriteRepo.Favorites["1-100"] != nil {
		t.Error("RemoveFavorite: favorite should be deleted")
	}
}

func TestComment_Success(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	result, err := svc.Comment(context.Background(), &CommentCommand{
		UserID:   1,
		VideoID:  100,
		ParentID: 0,
		Content:  "test comment",
	})
	if err != nil {
		t.Fatalf("Comment: expected nil, got %v", err)
	}
	if result.CommentID == 0 {
		t.Error("Comment: expected non-zero commentID")
	}
}

func TestComment_Reply(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[1] = &entity.Comment{
		ID:      1,
		UserID:  2,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	result, err := svc.Comment(context.Background(), &CommentCommand{
		UserID:   1,
		VideoID:  100,
		ParentID: 1,
		Content:  "reply",
	})
	if err != nil {
		t.Fatalf("Comment reply: expected nil, got %v", err)
	}
	if result.CommentID == 0 {
		t.Error("Comment reply: expected non-zero commentID")
	}
}

func TestComment_CreateError(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.CreateErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	_, err := svc.Comment(context.Background(), &CommentCommand{
		UserID:   1,
		VideoID:  100,
		ParentID: 0,
		Content:  "test",
	})
	if err == nil {
		t.Fatal("Comment: expected error, got nil")
	}
}

func TestDeleteComment_Success(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[1] = &entity.Comment{
		ID:      1,
		UserID:  1,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	err := svc.DeleteComment(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("DeleteComment: expected nil, got %v", err)
	}
	if commentRepo.Comments[1] != nil {
		t.Error("DeleteComment: comment should be deleted")
	}
}

func TestDeleteComment_NotOwner(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.Comments[1] = &entity.Comment{
		ID:      1,
		UserID:  2,
		VideoID: 100,
	}
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	err := svc.DeleteComment(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("DeleteComment: expected nil, got %v", err)
	}
	if commentRepo.Comments[1] == nil {
		t.Error("DeleteComment: comment should not be deleted when user is not owner")
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	svc := newTestService(mock.NewMockLikeRepository(), mock.NewMockCommentRepository(), mock.NewMockFavoriteRepository())

	err := svc.DeleteComment(context.Background(), 1, 999)
	if err != nil {
		t.Fatalf("DeleteComment: expected nil, got %v", err)
	}
}

func TestDeleteComment_RepositoryError(t *testing.T) {
	commentRepo := mock.NewMockCommentRepository()
	commentRepo.GetErr = errInternal
	svc := newTestService(mock.NewMockLikeRepository(), commentRepo, mock.NewMockFavoriteRepository())

	err := svc.DeleteComment(context.Background(), 1, 1)
	if err == nil {
		t.Fatal("DeleteComment: expected error, got nil")
	}
}
