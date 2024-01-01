package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHelloMessageHandler godoc
// @summary Health Check
// @description Health checking for the service
// @tags healthcheck
// @id GetHelloMessageHandler
// @produce json
// @param name query string false "name of the active user"
// @response 200 {object} handler.MessageResponse "OK"
// @router /hello [get]
func GetHelloMessage(c *gin.Context) {
	name := c.Query("name")

	if name != "" {
		c.JSON(makeHttpMessageResponse(http.StatusOK, fmt.Sprintf("Hello, %s!", name)))
		return
	}

	c.JSON(makeHttpMessageResponse(http.StatusOK, "Hello, Guest!"))
}
