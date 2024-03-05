package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		remoteArress := strings.Split(c.Request.RemoteAddr, ":")

		log.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.Path).
			Str("hostname", c.Request.Host).
			Str("ip", ip).
			Str("remote_address", remoteArress[0]).
			Str("remote_port", remoteArress[1]).
			Msg("incoming request")
		c.Next()
	}
}

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		c.Next()

		log.Info().
			Int("status_code", c.Writer.Status()).
			Dur("response_time", time.Since(now)).
			Msg("request completed")
	}
}
