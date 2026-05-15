package vectorstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type QdrantStore struct {
	addr           string
	collectionName string
}

type Document struct {
	ID       string
	Content  string
	Metadata map[string]any
	Score    float32
}

func (d *Document) GetScore() float32 {
	if d.Score != 0 {
		return d.Score
	}
	if s, ok := d.Metadata["score"].(float32); ok {
		return s
	}
	if s, ok := d.Metadata["score"].(float64); ok {
		return float32(s)
	}
	return 0
}

func NewQdrantStore(addr, collectionName string) (*QdrantStore, error) {
	return &QdrantStore{
		addr:           addr,
		collectionName: collectionName,
	}, nil
}

func (s *QdrantStore) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	url := fmt.Sprintf("%s%s", s.addr, path)
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *QdrantStore) Search(ctx context.Context, queryVector []float32, limit int) ([]*Document, error) {
	body := map[string]any{
		"vector":       queryVector,
		"limit":        limit,
		"with_payload": true,
	}

	path := fmt.Sprintf("/collections/%s/points/search", s.collectionName)
	data, err := s.doRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result []struct {
			ID      string         `json:"id"`
			Score   float64        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	docs := make([]*Document, len(result.Result))
	for i, r := range result.Result {
		docs[i] = &Document{
			ID:       r.ID,
			Score:    float32(r.Score),
			Content:  fmt.Sprintf("%v", r.Payload["content"]),
			Metadata: r.Payload,
		}
	}
	return docs, nil
}

func (s *QdrantStore) Upsert(ctx context.Context, id string, vector []float32, doc *Document) error {
	body := map[string]any{
		"points": []map[string]any{
			{
				"id":      id,
				"vector":  vector,
				"payload": doc.Metadata,
			},
		},
	}

	path := fmt.Sprintf("/collections/%s/points", s.collectionName)
	_, err := s.doRequest(ctx, "PUT", path, body)
	return err
}

func (s *QdrantStore) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/collections/%s/points/%s", s.collectionName, id)
	_, err := s.doRequest(ctx, "DELETE", path, nil)
	return err
}
