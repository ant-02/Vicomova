package query

import (
	"context"
	"fmt"
	"strconv"

	"vicomova/internal/search/domain/entity"
	"vicomova/internal/search/infrastructure/embedding"
	"vicomova/internal/search/infrastructure/vectorstore"
)

type SearchService struct {
	embedder  *embedding.OpenAIEmbedder
	store     *vectorstore.QdrantStore
	videoRepo VideoMetaRepository
}

type VideoMetaRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.VideoSearchResult, error)
}

func NewSearchService(
	embedder *embedding.OpenAIEmbedder,
	store *vectorstore.QdrantStore,
	videoRepo VideoMetaRepository,
) *SearchService {
	return &SearchService{
		embedder:  embedder,
		store:     store,
		videoRepo: videoRepo,
	}
}

func (s *SearchService) SearchVideos(ctx context.Context, q *SearchVideosQuery) (*SearchVideosResult, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if s.embedder == nil || s.store == nil {
		return &SearchVideosResult{
			Videos:     []*entity.VideoSearchResult{},
			NextCursor: "",
			HasMore:    false,
		}, nil
	}

	vector, err := s.embedder.EmbedString(ctx, q.Query)
	if err != nil {
		return nil, fmt.Errorf("embed query failed: %w", err)
	}

	docs, err := s.store.Search(ctx, vector, limit)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	if len(docs) == 0 {
		return &SearchVideosResult{
			Videos:     []*entity.VideoSearchResult{},
			NextCursor: "",
			HasMore:    false,
		}, nil
	}

	ids := make([]int64, 0, len(docs))
	scoreMap := make(map[int64]float32)
	for _, doc := range docs {
		idStr, ok := doc.Metadata["video_id"]
		if !ok {
			continue
		}
		var id int64
		switch v := idStr.(type) {
		case float64:
			id = int64(v)
		case int64:
			id = v
		case string:
			id, _ = strconv.ParseInt(v, 10, 64)
		}
		if id == 0 {
			continue
		}
		ids = append(ids, id)
		scoreMap[id] = doc.GetScore()
	}

	videos, err := s.videoRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("fetch video meta failed: %w", err)
	}

	for _, v := range videos {
		if score, ok := scoreMap[v.VideoID]; ok {
			v.Score = score
		}
	}

	return &SearchVideosResult{
		Videos:     videos,
		NextCursor: "",
		HasMore:    len(docs) >= limit,
	}, nil
}

func (s *SearchService) IndexVideo(ctx context.Context, q *IndexVideoQuery) error {
	if s.embedder == nil || s.store == nil {
		return nil
	}

	vector, err := s.embedder.EmbedString(ctx, q.Title+" "+q.Description)
	if err != nil {
		return fmt.Errorf("embed failed: %w", err)
	}

	doc := &vectorstore.Document{
		Metadata: map[string]any{
			"video_id": q.VideoID,
			"title":    q.Title,
			"content":  q.Title + " " + q.Description,
		},
	}

	return s.store.Upsert(ctx, fmt.Sprintf("%d", q.VideoID), vector, doc)
}

func (s *SearchService) DeleteVideoIndex(ctx context.Context, videoID int64) error {
	if s.store == nil {
		return nil
	}
	return s.store.Delete(ctx, fmt.Sprintf("%d", videoID))
}
