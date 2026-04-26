package gateway

import (
	interactionHdl "vicomova/internal/interaction/interfaces/http/handler"
	interactionRtr "vicomova/internal/interaction/interfaces/http/router"
	videoHdl "vicomova/internal/video/interfaces/http/handler"
	videoRtr "vicomova/internal/video/interfaces/http/router"
	"vicomova/internal/user/domain/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterVideoRoutes(h *server.Hertz, vh *videoHdl.VideoHandler, tokenSvc *service.TokenService) {
	videoRtr.RegisterRoutes(h, vh, tokenSvc)
}

func RegisterInteractionRoutes(h *server.Hertz, ih *interactionHdl.InteractionHandler) {
	interactionRtr.RegisterRoutes(h, ih)
}
