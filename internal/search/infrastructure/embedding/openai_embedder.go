package embedding

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Embedder interface {
	EmbedStrings(ctx context.Context, texts []string) ([][]float32, error)
	EmbedString(ctx context.Context, text string) ([]float32, error)
}

type OpenAIEmbedder struct {
	apiKey string
	apiURL string
	model  string
}

type Option func(*OpenAIEmbedder)

func WithModel(model string) Option {
	return func(e *OpenAIEmbedder) {
		e.model = model
	}
}

func WithAPIURL(apiURL string) Option {
	return func(e *OpenAIEmbedder) {
		e.apiURL = apiURL
	}
}

func NewOpenAIEmbedder(apiKey string, opts ...Option) (*OpenAIEmbedder, error) {
	e := &OpenAIEmbedder{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/embeddings",
		model:  "text-embedding-3-small",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e, nil
}

func (e *OpenAIEmbedder) EmbedStrings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqBody := map[string]any{
		"model": e.model,
		"input": texts,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.apiURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding API error: %s", string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(result.Data))
	for i, d := range result.Data {
		embeddings[i] = d.Embedding
	}
	return embeddings, nil
}

func (e *OpenAIEmbedder) EmbedString(ctx context.Context, text string) ([]float32, error) {
	embs, err := e.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embs) == 0 {
		return nil, nil
	}
	return embs[0], nil
}

var _ Embedder = (*OpenAIEmbedder)(nil)
