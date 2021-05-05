package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitUsersRoutes(groupApi *gin.RouterGroup) {
	auth := groupApi.Group("/auth")
	{
		auth.POST("/Sign-In", h.signIn)
		auth.POST("/Sign-Up", h.signUp)
		auth.POST("/refreshTokens", h.refreshTokens)

		api := auth.Group("/api", h.userIdentity)
		{
			ads := api.Group("/ads")
			{
				ads.GET("/", h.getAllAds)
				ads.POST("/", h.createAd)
				ads.GET("/:id", h.getAdById)
				ads.PUT("/:id", h.updateAd)
				ads.DELETE("/:id", h.deleteAd)
			}
		}
	}
}
