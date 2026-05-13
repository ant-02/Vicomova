package query

import (
	"vicomova/internal/video/domain/entity"
)

type GetVideoResult struct {
	Video *entity.Video
}

type ListVideoResult struct {
	Videos []*entity.Video
	Total  int64
}

type GetVideoStreamResult struct {
	Video *entity.Video
}

type GetVideoCoverResult struct {
	CoverURL string
}

type GetUploadTokenResult struct {
	Token  string
	Key    string
	Domain string
	Host   string // 七牛云上传地址
}

// HotVideoItem 热门视频项
type HotVideoItem struct {
	VideoID      int64
	Title        string
	CoverURL     string
	UserName     string
	Duration     int
	ViewCount    int64
	CommentCount int64
}

// ListHotVideosQuery 热门视频查询
type ListHotVideosQuery struct {
	Cursor string // 游标分页
	Limit  int
}

// ListHotVideosResult 热门视频结果
type ListHotVideosResult struct {
	Videos     []*HotVideoItem
	NextCursor string
	HasMore    bool
}

// PublishedVideosQuery 发布视频查询
type PublishedVideosQuery struct {
	Cursor string
	Limit  int
}

// PublishedVideosResult 发布视频结果
type PublishedVideosResult struct {
	Videos     []*HotVideoItem
	NextCursor string
	HasMore    bool
}

// CategoryVideosQuery 分类视频查询
type CategoryVideosQuery struct {
	Cursor     string
	Limit      int
	CategoryID int
}

// CategoryVideosResult 分类视频结果
type CategoryVideosResult struct {
	Videos     []*HotVideoItem
	NextCursor string
	HasMore    bool
}
