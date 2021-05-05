package v1

import (
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
)

func (h *Handler) signIn(ctx *gin.Context) {
	var input signInInput
	if err := ctx.BindJSON(&input); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.Authorization.SignIn(service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
		Ua:       ctx.GetHeader("User-Agent"),
		Ip:       ctx.Request.RemoteAddr,
	})
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

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

	ctx.JSON(http.StatusCreated, tokens)
}
