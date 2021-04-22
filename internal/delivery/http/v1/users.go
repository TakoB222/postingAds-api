package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitUsersRoutes(groupApi *gin.RouterGroup) {
	users := groupApi.Group("/users")
	{
		users.POST("/SignIn")
		users.POST("/SignUp")

		authenticated := users.Group("/", h.studentIdentity)
		{

		}
	}
}
