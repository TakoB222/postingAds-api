package v1

import (
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Services, tokeManager auth.TokenManager) *Handler {
	return &Handler{services: service, tokenManager: tokeManager}
}

func (h *Handler) Init(groupApi *gin.RouterGroup) {
	v1 := groupApi.Group("/v1")
	{
		h.InitUsersRoutes(v1)
	}
}
