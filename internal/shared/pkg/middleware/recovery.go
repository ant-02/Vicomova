package middleware

import (
	"context"
	"vicomova/internal/shared/pkg/log"

	"github.com/cloudwego/hertz/pkg/app"
)

func Recovery() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				log.Error.Printf("panic recovered: %v", err)
				ctx.JSON(500, map[string]interface{}{
					"code": 500,
					"msg":  "Internal server error",
				})
				ctx.Abort()
			}
		}()
		ctx.Next(c)
	}
}
