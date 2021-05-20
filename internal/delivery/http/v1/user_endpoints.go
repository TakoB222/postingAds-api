package v1

import (
	"errors"
	"fmt"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
			fts := api.Group("/fts")
			{
				fts.GET("/", h.fts)
			}
		}
	}
}

type (
	signInInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	signUpInput struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	refreshTokensInput struct {
		RefreshToken string `json:"RefreshToken" binding:"required"`
	}

	tokenResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	inputFTSRequest struct {
		Request string `json:"request" binding:"required"`
	}
)

//------------------Authorization------------------
// @Summary User SignIn
// @Tags users-auth
// @Description user sign in
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/Sign-In [post]
func (h *Handler) signIn(ctx *gin.Context) {
	var input signInInput
	if err := ctx.BindJSON(&input); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.Authorization.SignIn(service.SignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @Accept  json
// @Produce  json
// @Param input body signUpInput true "sign up info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/Sign-Up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var input signUpInput
	if err := ctx.BindJSON(&input); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	userId, err := h.services.Authorization.SignUp(service.UserSignUpInput{
		FirsName: input.FirstName,
		LastName: input.LastName,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": userId,
	})
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshTokensInput true "refresh token info"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/refreshTokens [post]
func (h *Handler) refreshTokens(ctx *gin.Context) {
	var refreshInput refreshTokensInput

	if err := ctx.BindJSON(&refreshInput); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.RefreshSession(service.RefreshInput{
		RefreshToken: refreshInput.RefreshToken,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

//------------------Ads------------------

//TODO: create input data struct validator

type (
	inputAd struct {
		Title       string        `json:"title" binding:"required"`
		Category    string        `json:"category" binding:"required"`
		Description string        `json:"description" binding:"required"`
		Price       int           `json:"price" binding:"required"`
		Contacts    inputContacts `json:"contacts" binding:"required"`
		Published   bool          `json:"published"`
		ImagesURL   []string      `json:"images_url" binding:"required"`
	}
	inputContacts struct {
		Name         string `json:"name" binding:"required"`
		Phone_number string `json:"phone_number" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Location     string `json:"location" binding:"required"`
	}
)

// @Summary User Get All His Ads
// @Security UsersAuth
// @Tags users-ads
// @Description user get all his ads by userId
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Ad
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/ads/ [get]
func (h *Handler) getAllAds(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ads, err := h.services.GetAllAds(userId)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ads)
}

// @Summary User Create Own Ad
// @Security UsersAuth
// @Tags users-ads
// @Description user create his own ad
// @Accept  json
// @Produce  json
// @Param input body inputAd true "create ad info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/ads/ [post]
func (h *Handler) createAd(ctx *gin.Context) {
	var inputAd inputAd
	if err := ctx.BindJSON(&inputAd); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if category := strings.Split(inputAd.Category, "/"); len(category) > 0 {
		inputAd.Category = category[len(category)-1]
	}

	fmt.Println(inputAd.Category)

	adId, err := h.services.Ad.CreateAd(userId, service.Ads{
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

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": adId,
	})
}

// @Summary User Get Ad By AdId
// @Security UsersAuth
// @Tags users-ads
// @Description user get ad by adId
// @Accept  json
// @Produce  json
// @Param id path string true "adId"
// @Success 200 {object} domain.Ad
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/ads/{id} [get]
func (h *Handler) getAdById(ctx *gin.Context) {
	adId := ctx.Param("id")

	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ad, err := h.services.GetAdById(userId, adId)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, ad)
}

// @Summary User Update His Ad
// @Security UsersAuth
// @Tags users-ads
// @Description user create ad
// @Accept  json
// @Produce  json
// @Param id path string true "adId"
// @Param input body inputAd true "ad info"
// @Success 200 {object} domain.Ad
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/ads/{id} [put]
func (h *Handler) updateAd(ctx *gin.Context) {
	adId := ctx.Param("id")

	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var inputAds inputAd
	if err = ctx.BindJSON(&inputAds); err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if category := strings.Split(inputAds.Category, "/"); len(category) > 0 {
		inputAds.Category = category[len(category)-1]
	}

	ad ,err := h.services.UpdateAd(userId, adId, service.Ads{
		Title:       inputAds.Title,
		Category:    inputAds.Category,
		Description: inputAds.Description,
		Price:       inputAds.Price,
		Contacts:    service.Contacts(inputAds.Contacts),
		Published:   inputAds.Published,
		ImagesURL:   inputAds.ImagesURL,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ad)
}

// @Summary User Delete Ad
// @Security UsersAuth
// @Tags users-ads
// @Description user delete his ad
// @Accept  json
// @Produce  json
// @Param id path string true "adId"
// @Success 200 {object} string "deleted"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/ads/{id} [delete]
func (h *Handler) deleteAd(ctx *gin.Context) {
	adId := ctx.Param("id")

	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.DeleteAd(userId, adId); err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "deleted")
}

// @Summary User Search Ads
// @Security UsersAuth
// @Tags users-ads
// @Description user search ads by his request string
// @Accept  json
// @Produce  json
// @Param input body inputFTSRequest true "search request"
// @Success 200 {object} repository.FtsResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/api/fts/ [get]
func (h *Handler) fts(ctx *gin.Context) {
	var input inputFTSRequest
	if err := ctx.BindJSON(&input); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	searchResult, err := h.services.Ad.Fts(input.Request)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, searchResult)
}

func getUserId(ctx *gin.Context) (string, error) {
	id, ok := ctx.Get(userContext)
	if !ok {
		return "", errors.New("empty user context")
	}

	if id == "" {
		return "", errors.New("empty body of user context")
	}

	userId, ok := id.(string)
	if !ok {
		return "", errors.New("invalid type of userId from context")
	}

	return userId, nil
}
