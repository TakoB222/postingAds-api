package v1

import (
	"github.com/TakoB222/postingAds-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(ctx *gin.Context, statusCode int, message string) {
	logger.Error(message)
	ctx.AbortWithStatusJSON(statusCode, response{message})
}
