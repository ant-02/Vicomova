package mysql

import (
	"vicomova/internal/video/domain/entity"
	videoVO "vicomova/internal/video/domain/valueobject"
)

func VideoToPO(v *entity.Video) *VideoPO {
	if v == nil {
		return nil
	}
	return &VideoPO{
		ID:           v.ID,
		UserID:       v.UserID,
		Title:        v.Title,
		Description:  v.Description,
		CoverURL:     v.CoverURL,
		VideoURL:     v.VideoURL,
		CategoryID:   v.CategoryID,
		ViewCount:    v.ViewCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		Duration:     v.Duration,
		HotScore:     v.HotScore,
		Status:       int8(v.Status),
	}
}

func POToVideo(po *VideoPO) *entity.Video {
	if po == nil {
		return nil
	}
	return &entity.Video{
		ID:           po.ID,
		UserID:       po.UserID,
		Title:        po.Title,
		Description:  po.Description,
		CoverURL:     po.CoverURL,
		VideoURL:     po.VideoURL,
		CategoryID:   po.CategoryID,
		ViewCount:    po.ViewCount,
		LikeCount:    po.LikeCount,
		CommentCount: po.CommentCount,
		Duration:     po.Duration,
		HotScore:     po.HotScore,
		Status:       videoVO.VideoStatus(po.Status),
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}
