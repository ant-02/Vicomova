package oss

import "fmt"

type OSSType string

const (
	OSSTypeQiniu OSSType = "qiniu"
)

// NewOSS 创建 OSS 实例（用于前端直传场景）
func NewOSS(ossType OSSType, cfg interface{}) (OSS, error) {
	switch ossType {
	case OSSTypeQiniu:
		qiniuCfg, ok := cfg.(*OSSConfig)
		if !ok {
			return nil, fmt.Errorf("invalid qiniu config")
		}
		return NewOSSQiniu(qiniuCfg), nil
	default:
		return nil, fmt.Errorf("unsupported oss type: %s", ossType)
	}
}
