package router

import (
	"vicomova/pkg/config"
	pkgMiddleware "vicomova/pkg/middleware"

	"github.com/cloudwego/hertz/pkg/app"
)

func protectedMw() []app.HandlerFunc {
	return []app.HandlerFunc{pkgMiddleware.Auth(config.Get().JWT.Secret)}
}
