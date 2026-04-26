package handlers

// ========== Common Types ==========

type ErrorResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// ========== Video Types ==========

type PublishVideoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  int32  `json:"category_id"`
	CoverUrl    string `json:"cover_url"`
	VideoUrl    string `json:"video_url"`
	Duration    int32  `json:"duration"`
}

type VideoStreamResponse struct {
	VideoURL string `json:"video_url"`
	Title    string `json:"title"`
}

type VideoItem struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CoverURL     string `json:"cover_url"`
	VideoURL     string `json:"video_url"`
	CategoryID   int32  `json:"category_id"`
	ViewCount    int64  `json:"view_count"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
	Duration     int32  `json:"duration"`
	CreatedAt    int64  `json:"created_at"`
}

type VideoListResponse struct {
	Videos []*VideoItem `json:"videos"`
	Total  int64         `json:"total"`
}

type CoverResponse struct {
	CoverURL string `json:"cover_url"`
}

// ========== Interaction Types ==========

type LikeRequest struct {
	VideoID int64 `json:"video_id"`
}

type UnlikeRequest struct {
	VideoID int64 `json:"video_id"`
}

type FavoriteRequest struct {
	VideoID int64 `json:"video_id"`
}

type UnfavoriteRequest struct {
	VideoID int64 `json:"video_id"`
}

type CommentRequest struct {
	VideoID  int64  `json:"video_id"`
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content"`
}

type DeleteCommentRequest struct {
	CommentID int64 `json:"comment_id"`
}

type LikeCommentRequest struct {
	CommentID int64 `json:"comment_id"`
}

type LikeItem struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	TargetType string `json:"target_type"`
	TargetID   int64  `json:"target_id"`
	CreatedAt  int64  `json:"created_at"`
}

type LikeListResponse struct {
	Likes []*LikeItem `json:"likes"`
	Total int64       `json:"total"`
}

type FavoriteItem struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	VideoID   int64  `json:"video_id"`
	CreatedAt int64  `json:"created_at"`
	Video     *VideoItem `json:"video,omitempty"`
}

type FavoriteListResponse struct {
	Favorites []*FavoriteItem `json:"favorites"`
	Total     int64           `json:"total"`
}

type CommentItem struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	VideoID   int64  `json:"video_id"`
	ParentID  int64  `json:"parent_id"`
	Content   string `json:"content"`
	LikeCount int64  `json:"like_count"`
	CreatedAt int64 `json:"created_at"`
}

type CommentListResponse struct {
	Comments []*CommentItem `json:"comments"`
	Total    int64          `json:"total"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
