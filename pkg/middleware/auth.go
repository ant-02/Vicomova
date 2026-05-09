package middleware

import (
	"context"
	"strings"

	"vicomova/pkg/constants"
	"vicomova/pkg/infrastructure/hertz"
	"vicomova/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

func Auth(jwtSecret string) app.HandlerFunc {
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

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			ctx.JSON(401, hertz.Fail(401, "Invalid token"))
			ctx.Abort()
			return
		}

		userID, err := utils.GetUserIDFromClaims(claims)
		if err != nil {
			ctx.JSON(401, hertz.Fail(401, "Invalid token claims"))
			ctx.Abort()
			return
		}

		ctx.Set(constants.ContextKeyUserID, userID)
		ctx.Next(c)
	}
}
