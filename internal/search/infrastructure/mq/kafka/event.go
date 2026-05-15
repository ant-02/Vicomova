package kafka

import (
	"encoding/json"
)

// VideoIndexEvent 视频索引事件
type VideoIndexEvent struct {
	VideoID     int64  `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (e *VideoIndexEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *VideoIndexEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}
