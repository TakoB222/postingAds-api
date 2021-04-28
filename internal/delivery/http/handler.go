package http

import (
	"github.com/TakoB222/postingAds-api/internal/delivery/http/v1"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, tokeManager auth.TokenManager) *Handler {
	return &Handler{services: service, tokenManager: tokeManager}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
