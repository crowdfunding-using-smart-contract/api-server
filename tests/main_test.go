package tests

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func NewTestContext(t *testing.T, recorder *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	ctx, router := gin.CreateTestContext(recorder)
	return ctx, router
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
