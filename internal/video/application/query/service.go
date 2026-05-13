package query

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	usergrpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/internal/video/domain/entity"
	"vicomova/internal/video/domain/repository"
	"vicomova/pkg/constants"
	"vicomova/pkg/infrastructure/oss"
	"vicomova/pkg/log"

	user "vicomova/third_party/kitex_gen/user"
)

type VideoQueryService struct {
	repo              repository.VideoRepository
	cache             repository.VideoCache
	hotCache          repository.HotVideoCache
	oss               oss.OSS
	viewCountProducer repository.ViewCountProducer
	userClient        *usergrpc.UserClient
	cacheRebuilding   atomic.Bool
}

func NewVideoQueryService(
	repo repository.VideoRepository,
	cache repository.VideoCache,
	ossClient oss.OSS,
	viewCountProducer repository.ViewCountProducer,
	hotCache repository.HotVideoCache,
	userClient *usergrpc.UserClient,
) *VideoQueryService {
	return &VideoQueryService{
		repo:              repo,
		cache:             cache,
		hotCache:          hotCache,
		oss:               ossClient,
		viewCountProducer: viewCountProducer,
		userClient:        userClient,
	}
}

// WarmUp 启动时预热，并启动后台定时刷新
func (s *VideoQueryService) WarmUp(ctx context.Context) error {
	log.Info.Printf("Hot cache warming up...")
	s.rebuildHotCache(ctx)

	// 启动后台定时刷新（只启动一次）
	hotCacheRefreshOnce.Do(func() {
		go s.runHotCacheRefresh()
		log.Info.Printf("Hot cache refresh background started")
	})

	return nil
}

var hotCacheRefreshOnce sync.Once

// runHotCacheRefresh 后台定时刷新热门视频缓存
func (s *VideoQueryService) runHotCacheRefresh() {
	ticker := time.NewTicker(constants.HotCacheRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.rebuildHotCache(context.Background())
		}
	}
}

// rebuildHotCache 异步重建热门视频缓存（原子性 swap）
func (s *VideoQueryService) rebuildHotCache(ctx context.Context) {
	log.Debug.Printf("rebuildHotCache: started")

	videos, err := s.repo.ListHot(ctx, 200) // 直接用 repo 的 ListHot（利用索引）
	if err != nil {
		log.Error.Printf("rebuildHotCache: ListHot failed: %v", err)
		return
	}
	log.Debug.Printf("rebuildHotCache: ListHot returned %d videos", len(videos))

	if len(videos) == 0 {
		log.Debug.Printf("rebuildHotCache: no videos to cache, returning")
		return
	}

	// 1. 写入 staging key
	log.Debug.Printf("rebuildHotCache: writing to staging key...")
	if err := s.hotCache.SetHotVideoIDsStaging(ctx, videos); err != nil {
		log.Error.Printf("rebuildHotCache: SetHotVideoIDsStaging failed: %v", err)
		return
	}
	log.Debug.Printf("rebuildHotCache: staging write done")

	// 2. 原子性 swap（瞬间切换，读者看到完整数据）
	log.Debug.Printf("rebuildHotCache: swapping staging to main key...")
	if err := s.hotCache.SwapHotVideoIDs(ctx); err != nil {
		log.Error.Printf("rebuildHotCache: SwapHotVideoIDs failed: %v", err)
		return
	}
	log.Debug.Printf("rebuildHotCache: swap done")

	log.Info.Printf("rebuildHotCache: refreshed %d videos", len(videos))
}

// buildHotVideoItems 将视频列表转换为 HotVideoItem（用于降级场景）
func (s *VideoQueryService) buildHotVideoItems(ctx context.Context, videos []*entity.Video, offset int64) []*HotVideoItem {
	if len(videos) == 0 {
		return []*HotVideoItem{}
	}

	// 收集用户 ID
	userIDSet := make(map[int64]struct{})
	for _, v := range videos {
		userIDSet[v.UserID] = struct{}{}
	}
	userIDs := make([]int64, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	// 批量获取用户信息
	userMap := make(map[int64]*user.User)
	if s.userClient != nil && len(userIDs) > 0 {
		users, err := s.userClient.BatchGetUsers(ctx, userIDs)
		if err != nil {
			log.Error.Printf("buildHotVideoItems: BatchGetUsers failed: %v", err)
		} else {
			userMap = users
		}
	}

	// 转换
	items := make([]*HotVideoItem, 0, len(videos))
	for _, v := range videos {
		item := &HotVideoItem{
			VideoID:      v.ID,
			Title:        v.Title,
			CoverURL:     v.CoverURL,
			Duration:     v.Duration,
			ViewCount:    v.ViewCount,
			CommentCount: v.CommentCount,
		}
		if u, ok := userMap[v.UserID]; ok {
			item.UserName = u.Username
		}
		items = append(items, item)
	}

	return items
}
