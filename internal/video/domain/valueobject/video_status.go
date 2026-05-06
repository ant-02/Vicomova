package valueobject

type VideoStatus int8

const (
	VideoStatusPending   VideoStatus = 0
	VideoStatusPublished VideoStatus = 1
	VideoStatusRemoved   VideoStatus = 2
)
