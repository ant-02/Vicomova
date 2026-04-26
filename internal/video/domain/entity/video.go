package entity

import "time"

type Video struct {
    ID           int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CoverURL    string    `json:"cover_url"`
    VideoURL    string    `json:"video_url"`
    CategoryID  int       `json:"category_id"`
    ViewCount   int64     `json:"view_count"`
    LikeCount   int64     `json:"like_count"`
    CommentCount int64     `json:"comment_count"`
    Duration    int       `json:"duration"`
    Status      int8      `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

const (
    VideoStatusPending  int8 = 0 // 审核中
    VideoStatusPublished int8 = 1 // 已发布
    VideoStatusRemoved   int8 = 2 // 下架
)

func (v *Video) IsPublished() bool {
    return v.Status == VideoStatusPublished
}
