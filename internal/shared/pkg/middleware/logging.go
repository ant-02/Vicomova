package middleware

import (
	"context"
	"time"
	"vicomova/internal/shared/pkg/log"

	"github.com/cloudwego/hertz/pkg/app"
)

func Logging() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.Path())
		method := string(ctx.Request.Method())

		ctx.Next(c)

		log.Info.Printf("method=%s path=%s status=%d duration=%v",
			method, path, ctx.Response.StatusCode(), time.Since(start))
	}
}