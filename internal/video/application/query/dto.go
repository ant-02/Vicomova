package query

import (
	"vicomova/internal/video/domain/entity"
)

type GetVideoResult struct {
	Video *entity.Video
}

type ListVideoResult struct {
	Videos []*entity.Video
	Total   int64
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
