package repository

import "context"

// ViewCountProducer 播放量消息生产者接口
type ViewCountProducer interface {
	// Record 发送单条播放量增量
	Record(ctx context.Context, videoID int64) error
	// RecordBatch 批量发送播放量增量
	RecordBatch(ctx context.Context, videoIDs []int64) error
}
