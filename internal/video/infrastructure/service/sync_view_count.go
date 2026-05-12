package service

import (
	"context"
	"sync"
	"time"

	"vicomova/internal/video/domain/repository"
	redisCache "vicomova/internal/video/infrastructure/persistence/redis"
	"vicomova/pkg/log"
)

var syncOnce sync.Once

// StartViewCountSync 启动后台播放量同步任务（只启动一次）
func StartViewCountSync(repo repository.VideoRepository) {
	syncOnce.Do(func() {
		go runViewCountSync(repo)
		log.Info.Printf("ViewCountSync: started")
	})
}

func runViewCountSync(repo repository.VideoRepository) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := syncViewCounts(context.Background(), repo); err != nil {
			log.Error.Printf("ViewCountSync: sync failed: %v", err)
		}
	}
}

func syncViewCounts(ctx context.Context, repo repository.VideoRepository) error {
	// 扫描所有待同步的播放量
	counts, err := redisCache.GetPendingViewCounts(ctx)
	if err != nil {
		return err
	}

	if len(counts) == 0 {
		return nil
	}

	log.Info.Printf("ViewCountSync: syncing %d videos", len(counts))

	// 批量更新到 MySQL
	for videoID, count := range counts {
		if err := repo.IncrementViewBatch(ctx, videoID, count); err != nil {
			log.Error.Printf("ViewCountSync: failed to sync video %d: %v", videoID, err)
			continue
		}
		// 删除 Redis key
		if err := redisCache.DelViewCountKey(ctx, videoID); err != nil {
			log.Error.Printf("ViewCountSync: failed to del key for video %d: %v", videoID, err)
		}
	}

	return nil
}
