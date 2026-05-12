package kafka

import "time"

const (
	// 窗口配置
	FlushInterval  = 30 * time.Second
	FlushThreshold = 100
)
