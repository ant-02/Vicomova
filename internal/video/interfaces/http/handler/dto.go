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

type UploadTokenRequest struct {
	VideoID    int64 `json:"video_id" form:"video_id"`
	UploadType int32 `json:"upload_type" form:"upload_type"` // 1=video, 2=cover
}

type UploadTokenResponse struct {
	Token  string `json:"token"`
	Key    string `json:"key"`
	Domain string `json:"domain"`
	Host   string `json:"host"` // 七牛云上传地址
}
