package v1

import (
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) InitAdminRoutes(groupApi *gin.RouterGroup) {
	admins := groupApi.Group("/admins")
	{
		admins.POST("/Sign-In", h.adminSignIn)
		admins.POST("/refreshTokens", h.adminRefreshTokens)

		api := admins.Group("/api", h.adminIdentity)
		{
			ads := api.Group("/ads")
			{
				ads.GET("/", h.adminGetAllAds)
				ads.GET("/:id", h.adminGetAd)
				ads.PUT("/:id", h.adminUpdateAd)
				ads.DELETE("/:id", h.adminDeleteAd)
			}
		}
	}
}

type (
	adminSignInInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	adminUpdateAdInput struct {
		Title       string                   `json:"title" binding:"required"`
		Category    string                   `json:"category" binding:"required"`
		Description string                   `json:"description" binding:"required"`
		Price       int                      `json:"price" binding:"required"`
		Contacts    adminInputUpdateContacts `json:"contacts" binding:"required"`
		Published   bool                     `json:"published"`
		ImagesURL   []string                 `json:"images_url" binding:"required"`
	}
	adminInputUpdateContacts struct {
		Name         string `json:"name" binding:"required"`
		Phone_number string `json:"phone_number" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Location     string `json:"location" binding:"required"`
	}
)

func (h *Handler) adminSignIn(ctx *gin.Context) {
	var input adminSignInInput
	if err := ctx.BindJSON(&input); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.Admin.AdminSignIn(service.SignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *Handler) adminRefreshTokens(ctx *gin.Context) {
	var refreshInput refreshTokensInput

	if err := ctx.BindJSON(&refreshInput); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.AdminRefreshSession(service.RefreshInput{
		RefreshToken: refreshInput.RefreshToken,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *Handler) adminGetAllAds(ctx *gin.Context) {
	ads, err := h.services.Admin.AdminGetAllAdsByAdmin()
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ads)
}

func (h *Handler) adminGetAd(ctx *gin.Context) {
	id := ctx.Param("id")

	ad, err := h.services.Admin.AdminGetAd(id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ad)
}

func (h *Handler) adminDeleteAd(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.services.Admin.AdminDeleteUserAdById(id); err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "deleted")
}

func (h *Handler) adminUpdateAd(ctx *gin.Context) {
	adId := ctx.Param("id")

	var inputAd adminUpdateAdInput
	if err := ctx.BindJSON(&inputAd); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	if category := strings.Split(inputAd.Category, "/"); len(category) > 0 {
		inputAd.Category = category[len(category)-1]
	}

	err := h.services.AdminUpdateAd(adId, service.Ads{
		Title:       inputAd.Title,
		Category:    inputAd.Category,
		Description: inputAd.Description,
		Price:       inputAd.Price,
		Contacts:    service.Contacts(inputAd.Contacts),
		Published:   inputAd.Published,
		ImagesURL:   inputAd.ImagesURL,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "updated")
}
