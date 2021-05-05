package v1

import (
	"errors"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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

func (h *Handler) getAllAds(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ads, err := h.services.GetAllAds(userId)
	if err != nil {
		if err.Error() == repository.ErrorEmptyResult {
			newResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusFound, ads)
}

func (h *Handler) createAd(ctx *gin.Context) {
	var inputAds inputAd
	if err := ctx.BindJSON(&inputAds); err != nil {
		newResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if category := strings.Split(inputAds.Category, "/"); len(category) > 0 {
		inputAds.Category = category[len(category)-1]
	} else {
		newResponse(ctx, http.StatusBadRequest, "empty category body")
		return
	}

	adId, err := h.services.Ad.CreateAd(userId, service.Ads{
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

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": adId,
	})
}

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
	} else {
		newResponse(ctx, http.StatusBadRequest, "empty category body")
		return
	}

	err = h.services.UpdateAd(userId, adId, service.Ads{
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

	ctx.JSON(http.StatusOK, "updated")
}

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

func getUserId(ctx *gin.Context) (string, error) {
	id, ok := ctx.Get(userContext)
	if !ok {
		return "", errors.New("empty user context")
	}

	userId, ok := id.(string)
	if !ok {
		return "", errors.New("invalid type of userId from context")
	}

	return userId, nil
}