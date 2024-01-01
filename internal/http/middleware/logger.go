package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	logger "github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.WithFields(logger.Fields{
			"method":        c.Request.Method,
			"url":           c.Request.URL.Path,
			"hostname":      c.Request.Host,
			"remoteAddress": strings.Split(c.Request.RemoteAddr, ":")[0],
			"remotePort":    strings.Split(c.Request.RemoteAddr, ":")[1],
		}).Infof("incoming request")

		c.Next()
	}
}

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		c.Next()

		logger.WithFields(logger.Fields{
			"statusCode":   c.Writer.Status(),
			"responseTime": time.Since(now),
		}).Info("reqeust completed")
	}
}
