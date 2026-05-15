package query

import "vicomova/internal/search/domain/entity"

type SearchVideosQuery struct {
	Query  string
	Limit  int
	Cursor string
}

type SearchVideosResult struct {
	Videos     []*entity.VideoSearchResult
	NextCursor string
	HasMore    bool
}

type IndexVideoQuery struct {
	VideoID     int64
	Title       string
	Description string
}
