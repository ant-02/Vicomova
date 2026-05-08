package middleware

import (
	"context"
	"strings"

	"vicomova/internal/user/domain/service"
	"vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(tokenSvc *service.TokenService) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		authHeader := string(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(401, hertz.Fail(401, "Missing authorization header"))
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(401, hertz.Fail(401, "Invalid authorization format"))
			ctx.Abort()
			return
		}

		claims, err := tokenSvc.ParseAccessToken(parts[1])
		if err != nil {
			ctx.JSON(401, hertz.Fail(401, "Invalid token"))
			ctx.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			ctx.JSON(401, hertz.Fail(401, "Invalid token claims"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", int64(userID))
		ctx.Next(c)
	}
}
