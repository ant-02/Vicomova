package query

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vicomova/internal/video/domain/entity"
	"vicomova/pkg/constants"
	videoErr "vicomova/pkg/errors"
	"vicomova/pkg/log"

	user "vicomova/third_party/kitex_gen/user"
)

// GetVideoStream 获取视频流信息（同时触发播放量统计）
func (s *VideoQueryService) GetVideoStream(ctx context.Context, videoID int64) (*GetVideoStreamResult, error) {
	// 先查缓存，命中则刷新TTL
	var video *entity.Video
	if s.cache != nil {
		cached, err := s.cache.GetAndRefresh(ctx, videoID)
		if err == nil && cached != nil {
			if !cached.IsPublished() {
				return nil, videoErr.ErrVideoNotPublished
			}
			video = cached
		}
	}

	// 缓存未命中，查数据库
	if video == nil {
		var err error
		video, err = s.repo.GetByID(ctx, videoID)
		if err != nil {
			return nil, err
		}
		if video == nil {
			return nil, videoErr.ErrVideoNotFound
		}
		if !video.IsPublished() {
			return nil, videoErr.ErrVideoNotPublished
		}

		// 写入缓存
		if s.cache != nil {
			s.cache.Set(ctx, video)
		}
	}

	// 触发播放量统计
	if s.viewCountProducer != nil {
		s.viewCountProducer.Record(ctx, videoID)
	}

	return &GetVideoStreamResult{Video: video}, nil
}

// ListByCategory 按分类列出视频
func (s *VideoQueryService) ListByCategory(ctx context.Context, categoryID, page, size int) (*ListVideoResult, error) {
	videos, total, err := s.repo.ListByCategory(ctx, categoryID, page, size)
	if err != nil {
		return nil, err
	}
	return &ListVideoResult{Videos: videos, Total: total}, nil
}

// ListByUser 列出用户发布的视频
func (s *VideoQueryService) ListByUser(ctx context.Context, userID int64, page, size int) (*ListVideoResult, error) {
	videos, total, err := s.repo.ListByUser(ctx, userID, page, size)
	if err != nil {
		return nil, err
	}
	return &ListVideoResult{Videos: videos, Total: total}, nil
}

// ListPublishedVideos 获取用户发布的视频（游标分页）
func (s *VideoQueryService) ListPublishedVideos(ctx context.Context, userID int64, q *PublishedVideosQuery) (*PublishedVideosResult, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := int64(0)
	if q.Cursor != "" {
		var err error
		offset, err = strconv.ParseInt(q.Cursor, 10, 64)
		if err != nil || offset < 0 {
			offset = 0
		}
	}

	videos, total, err := s.repo.ListByUser(ctx, userID, int(offset/int64(limit))+1, limit)
	if err != nil {
		return nil, err
	}

	items := make([]*HotVideoItem, 0, len(videos))
	for _, v := range videos {
		items = append(items, &HotVideoItem{
			VideoID:      v.ID,
			Title:        v.Title,
			CoverURL:     v.CoverURL,
			Duration:     v.Duration,
			ViewCount:    v.ViewCount,
			CommentCount: v.CommentCount,
		})
	}

	nextOffset := offset + int64(len(videos))
	var nextCursor string
	hasMore := nextOffset < total
	if hasMore {
		nextCursor = strconv.FormatInt(nextOffset, 10)
	}

	return &PublishedVideosResult{
		Videos:     items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// GetUploadToken 获取视频封面，供前端直传到 OSS（video_id 生成唯一 key）
// key 格式：{video_id}/{upload_type}/{year}/{month}/{day}/{timestamp}
// uploadType: 1=video, 2=cover
func (s *VideoQueryService) GetUploadToken(ctx context.Context, videoID int64, uploadType int32) (*GetUploadTokenResult, error) {
	// 检查视频是否存在
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	if video == nil {
		return nil, fmt.Errorf("video not found")
	}

	// 生成 key：{video_id}/{upload_type}/{year}/{month}/{day}/{timestamp}
	now := time.Now()
	uploadTypeStr := map[int32]string{constants.UploadTokenTypeVideo: "video", constants.UploadTokenTypeCover: "cover"}[uploadType]
	key := fmt.Sprintf("%d/%s/%d/%02d/%02d/%d", videoID, uploadTypeStr, now.Year(), now.Month(), now.Day(), now.Unix())

	expire := constants.UploadTokenExpire * time.Second

	token, host, domain, err := s.oss.GetUploadToken(ctx, key, expire)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload token: %w", err)
	}

	return &GetUploadTokenResult{
		Token:  token,
		Key:    key,
		Domain: domain,
		Host:   host,
	}, nil
}

// ListHotVideos 列出热门视频（分页）
func (s *VideoQueryService) ListHotVideos(ctx context.Context, q *ListHotVideosQuery) (*ListHotVideosResult, error) {
	// 参数校验
	limit := q.Limit
	if limit <= 0 {
		limit = constants.DefaultHotVideoLimit
	}
	if limit > constants.MaxHotVideoLimit {
		limit = constants.MaxHotVideoLimit
	}

	// 解析游标
	var offset int64 = 0
	if q.Cursor != "" {
		var err error
		offset, err = strconv.ParseInt(q.Cursor, 10, 64)
		if err != nil || offset < 0 {
			offset = 0
		}
	}

	// 1. 从缓存获取热门视频 ID 列表
	scores, err := s.hotCache.GetHotVideoIDs(ctx, offset, int64(limit))
	if err != nil {
		log.Error.Printf("ListHotVideos: GetHotVideoIDs failed: %v", err)
	}

	// 2. 缓存未命中，触发异步回填并降级返回
	if len(scores) == 0 {
		// 双重检查锁：只允许一个请求触发回填
		if s.cacheRebuilding.CompareAndSwap(false, true) {
			go func() {
				defer s.cacheRebuilding.Store(false)
				s.rebuildHotCache(context.Background())
			}()
		}

		// 降级：直接用 repo.ListHot 返回
		fallbackVideos, err := s.repo.ListHot(ctx, limit)
		if err != nil {
			log.Error.Printf("ListHotVideos: hotAlgo.GetHotVideos fallback failed: %v", err)
			return &ListHotVideosResult{
				Videos:     []*HotVideoItem{},
				NextCursor: "",
				HasMore:    false,
			}, nil
		}

		items := s.buildHotVideoItems(ctx, fallbackVideos, offset)
		return &ListHotVideosResult{
			Videos:     items,
			NextCursor: "",
			HasMore:    false,
		}, nil
	}

	// 3. 提取视频 ID
	videoIDs := make([]int64, len(scores))
	for i, score := range scores {
		videoIDs[i] = score.VideoID
	}

	// 4. 批量获取视频元数据
	metas, missIDs, err := s.getVideoMetas(ctx, videoIDs)
	if err != nil {
		log.Error.Printf("ListHotVideos: getVideoMetas failed: %v", err)
		return nil, err
	}

	// 5. 缓存未命中时从数据库加载并缓存
	if len(missIDs) > 0 {
		for _, id := range missIDs {
			v, err := s.repo.GetByID(ctx, id)
			if err != nil || v == nil {
				continue
			}
			if !v.IsPublished() {
				continue
			}
			meta := &entity.HotVideoMeta{}
			meta.FromVideo(v)
			metas[id] = meta
			// 异步缓存
			go func(videoID int64, m *entity.HotVideoMeta) {
				s.hotCache.SetVideoMeta(context.Background(), videoID, m)
			}(id, meta)
		}
	}

	// 6. 收集所有需要的用户 ID
	userIDSet := make(map[int64]struct{})
	for _, meta := range metas {
		if meta != nil {
			userIDSet[meta.UserID] = struct{}{}
		}
	}
	userIDs := make([]int64, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	// 7. 批量获取用户信息
	var userMap map[int64]*user.User
	if s.userClient != nil && len(userIDs) > 0 {
		users, err := s.userClient.BatchGetUsers(ctx, userIDs)
		if err != nil {
			log.Error.Printf("ListHotVideos: BatchGetUsers failed: %v", err)
		} else {
			userMap = users
		}
	}

	// 8. 组装结果
	items := make([]*HotVideoItem, 0, len(videoIDs))
	for _, videoID := range videoIDs {
		meta := metas[videoID]
		if meta == nil {
			continue
		}
		item := &HotVideoItem{
			VideoID:      meta.VideoID,
			Title:        meta.Title,
			CoverURL:     meta.CoverURL,
			Duration:     meta.Duration,
			ViewCount:    meta.ViewCount,
			CommentCount: meta.CommentCount,
		}
		if userMap != nil {
			if u, ok := userMap[meta.UserID]; ok {
				item.UserName = u.Username
			}
		}
		items = append(items, item)
	}

	// 9. 计算下一页游标
	total, err := s.hotCache.GetHotVideoCount(ctx)
	if err != nil {
		log.Error.Printf("ListHotVideos: GetHotVideoCount failed: %v", err)
	}

	nextOffset := offset + int64(len(videoIDs))
	var nextCursor string
	hasMore := nextOffset < total
	if hasMore {
		nextCursor = strconv.FormatInt(nextOffset, 10)
	}

	return &ListHotVideosResult{
		Videos:     items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// getVideoMetas 批量获取视频元数据
func (s *VideoQueryService) getVideoMetas(ctx context.Context, videoIDs []int64) (map[int64]*entity.HotVideoMeta, []int64, error) {
	metas := make(map[int64]*entity.HotVideoMeta)
	var missIDs []int64

	for _, id := range videoIDs {
		if s.hotCache == nil {
			continue
		}
		meta, err := s.hotCache.GetVideoMeta(ctx, id)
		if err != nil || meta == nil {
			missIDs = append(missIDs, id)
		} else {
			metas[id] = meta
		}
	}

	return metas, missIDs, nil
}
