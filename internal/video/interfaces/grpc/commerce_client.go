package grpc

import (
	"context"

	commerce "vicomova/third_party/kitex_gen/commerce"
	commerceservice "vicomova/third_party/kitex_gen/commerce/commerceservice"

	"github.com/cloudwego/kitex/client"
)

type CommerceClient struct {
	cli commerceservice.Client
}

func NewCommerceClient(serviceName, addr string) (*CommerceClient, error) {
	cli, err := commerceservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &CommerceClient{cli: cli}, nil
}

func (c *CommerceClient) CheckVideoAccess(ctx context.Context, userID int64, videoID int64) (*commerce.CheckVideoAccessResponse, error) {
	return c.cli.CheckVideoAccess(ctx, &commerce.CheckVideoAccessRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}
