package handler

type ErrorResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type SaveVideoRequest struct {
	VideoID     int64  `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  int32  `json:"category_id"`
	CoverUrl    string `json:"cover_url"`
	VideoUrl    string `json:"video_url"`
	Duration    int32  `json:"duration"`
}

type SubmitVideoRequest struct {
	VideoID int64 `json:"video_id"`
}

type PublishVideoRequest struct {
	VideoID     int64  `json:"video_id"`
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
	Total  int64        `json:"total"`
}

type CoverResponse struct {
	CoverURL string `json:"cover_url"`
}

type UploadTokenRequest struct {
	Key          string `json:"key" binding:"required"` // 文件名
	ExpireSeconds int64  `json:"expire_seconds"`        // token 过期时间（秒），默认 3600
}

type UploadTokenResponse struct {
	Token  string `json:"token"`
	Key    string `json:"key"`
	Domain string `json:"domain"`
}
