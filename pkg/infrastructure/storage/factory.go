package storage

import "fmt"

type StorageType string

const (
	StorageTypeLocal  StorageType = "local"
	StorageTypeQiniu StorageType = "qiniu"
)

func NewImageStorage(storageType StorageType, cfg interface{}) (ImageStorage, error) {
	switch storageType {
	case StorageTypeLocal:
		localCfg, ok := cfg.(*LocalImageConfig)
		if !ok {
			return nil, fmt.Errorf("invalid local config")
		}
		return NewLocalImageStorage(localCfg), nil
	case StorageTypeQiniu:
		qiniuCfg, ok := cfg.(*QiniuConfig)
		if !ok {
			return nil, fmt.Errorf("invalid qiniu config")
		}
		return NewQiniuStorage(qiniuCfg), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
