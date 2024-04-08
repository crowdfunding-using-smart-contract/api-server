package middleware

import (
	"errors"
	"fmt"
	"fund-o/api-server/pkg/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := parseQueryToken(c, tokenMaker)
		if err == nil {
			c.Set(AuthorizationPayloadKey, payload)
			c.Next()
			return
		}

		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err = tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}

func parseQueryToken(c *gin.Context, tokenMaker token.Maker) (*token.Payload, error) {
	accessToken := c.Query("token")
	if len(accessToken) == 0 {
		return nil, errors.New("access token is not provided")
	}

	payload, err := tokenMaker.VerifyToken(accessToken)
	return payload, err
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status":      http.StatusText(http.StatusUnauthorized),
		"status_code": http.StatusUnauthorized,
		"error":       err.Error(),
	}
}
