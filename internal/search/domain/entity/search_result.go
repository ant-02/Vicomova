package entity

type VideoSearchResult struct {
	VideoID      int64
	Title        string
	Description  string
	CoverURL     string
	ViewCount    int64
	LikeCount    int64
	CommentCount int64
	Score        float32
}
