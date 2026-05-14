package handler

type ErrorResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

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
	Likes      []*LikeItem `json:"likes"`
	NextCursor string      `json:"next_cursor"`
	HasMore    bool        `json:"has_more"`
}

type FavoriteItem struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	VideoID   int64 `json:"video_id"`
	CreatedAt int64 `json:"created_at"`
}

type FavoriteListResponse struct {
	Favorites  []*FavoriteItem `json:"favorites"`
	NextCursor string          `json:"next_cursor"`
	HasMore    bool            `json:"has_more"`
}

type CommentItem struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	VideoID   int64  `json:"video_id"`
	ParentID  int64  `json:"parent_id"`
	Content   string `json:"content"`
	LikeCount int64  `json:"like_count"`
	CreatedAt int64  `json:"created_at"`
}

type CommentListResponse struct {
	Comments   []*CommentItem `json:"comments"`
	NextCursor string         `json:"next_cursor"`
	HasMore    bool           `json:"has_more"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
