package valueobject

type VideoStatus int8

const (
	VideoStatusEditing   VideoStatus = 0 // 编辑中
	VideoStatusPending   VideoStatus = 1 // 审核中
	VideoStatusPublished VideoStatus = 2 // 已发布
	VideoStatusRemoved   VideoStatus = 3 // 下架
)
