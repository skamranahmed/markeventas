package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/markeventas/internal/token"
	"github.com/skamranahmed/markeventas/pkg/log"
)

const (
	AuthorizationHeaderKey   = "Authorization"
	AauthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey  = "AuthorizationPayload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			errMsg := fmt.Sprintf("authorization header is required")
			log.Error(errMsg)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			errMsg := fmt.Sprintf("invalid authorization header format")
			log.Error(errMsg)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AauthorizationTypeBearer {
			errMsg := fmt.Sprintf("incorrect authorization type")
			log.Error(errMsg)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			errMsg := fmt.Sprintf("token verification failed, error: %v", err)
			log.Error(errMsg)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}
