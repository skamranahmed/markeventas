package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
)

type UserHandler interface {
	TwitterOAuthLogin(c *gin.Context)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		service: userService,
	}
}

type userHandler struct {
	service service.UserService
}

func (uh *userHandler) TwitterOAuthLogin(c *gin.Context) {
	var requestPayload UserTwitterOAuthLoginInput
	err := c.BindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}
}
