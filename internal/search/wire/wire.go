package wire

import (
	"vicomova/internal/search/application/query"
	"vicomova/internal/search/infrastructure/embedding"
	infraMysql "vicomova/internal/search/infrastructure/persistence/mysql"
	"vicomova/internal/search/infrastructure/vectorstore"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	SearchService *query.SearchService
}

func NewProvider() (*Provider, error) {
	cfg := config.Get()

	if err := mysql.Init(&cfg.Database); err != nil {
		return nil, err
	}

	if err := redis.Init(&cfg.Redis); err != nil {
		return nil, err
	}

	mysqlClient := mysql.GetClient()

	videoMetaRepo := infraMysql.NewVideoMetaRepository(mysqlClient)

	var embedder *embedding.OpenAIEmbedder
	var store *vectorstore.QdrantStore

	openAIKey := cfg.OpenAI.APIKey
	openAIURL := cfg.OpenAI.APIURL
	openAIModel := cfg.OpenAI.Model
	qdrantAddr := cfg.Qdrant.Addr
	collectionName := "videos"

	if openAIKey != "" && qdrantAddr != "" {
		var err error
		opts := []embedding.Option{}
		if openAIURL != "" {
			opts = append(opts, embedding.WithAPIURL(openAIURL))
		}
		if openAIModel != "" {
			opts = append(opts, embedding.WithModel(openAIModel))
		}
		embedder, err = embedding.NewOpenAIEmbedder(openAIKey, opts...)
		if err != nil {
			log.Warn.Printf("failed to create embedder: %v", err)
		}
		store, err = vectorstore.NewQdrantStore(qdrantAddr, collectionName)
		if err != nil {
			log.Warn.Printf("failed to create qdrant store: %v", err)
		}
	}

	searchSvc := query.NewSearchService(embedder, store, videoMetaRepo)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
	})

	return &Provider{
		SearchService: searchSvc,
	}, nil
}

type ServiceKey string

const (
	ServiceSearch ServiceKey = constants.ServiceSearch
)
