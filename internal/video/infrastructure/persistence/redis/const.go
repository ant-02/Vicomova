package redis

import "time"

const (
	VideoCacheKeyPrefix = "video:"
	ViewCountKeyPrefix  = "video:%d:views" // 播放量增量 key
	VideoMetaTTL        = time.Hour        // 视频元数据缓存 1h
)
