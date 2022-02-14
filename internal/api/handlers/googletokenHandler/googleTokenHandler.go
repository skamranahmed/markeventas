package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/middlewares"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

type GoogleTokenHandler interface {
	SaveRefreshToken(c *gin.Context)
}

func NewGoogleTokenHandler() GoogleTokenHandler {
	return &googleTokenHandler{}
}

type googleTokenHandler struct{}

type RefreshTokenRequest struct {
	// we use this code to fetch the refresh token
	// works only for the first time per-account
	Code string `json:"code" binding:"required"`
}

func (gh *googleTokenHandler) SaveRefreshToken(c *gin.Context) {
	// extract the payload from the context that was set by the AuthMiddleware
	jwtAuthPayload, exists := c.Get(middlewares.AuthorizationPayloadKey)
	if !exists {
		log.Errorf("authorization payload is not present")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtAuthToken, ok := jwtAuthPayload.(*token.Payload)
	if !ok {
		log.Errorf("invalid authorization payload, got payload: %+v", jwtAuthToken)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var requestPayload RefreshTokenRequest
	err := c.ShouldBindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	return
}
