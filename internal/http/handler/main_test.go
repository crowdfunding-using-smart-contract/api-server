package handler

import (
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/pkg/password"
	"fund-o/api-server/pkg/random"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	userID string,
	duration time.Duration,
) {
	tkn, payload, err := tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tkn)
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}

func randomUser(t *testing.T) entity.User {
	hashedPassword, err := password.HashPassword("@Password123")
	require.NoError(t, err)

	birthDate, err := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	require.NoError(t, err)

	user := entity.User{
		Base: entity.Base{
			ID: uuid.New(),
		},
		Email:          random.NewEmail(),
		HashedPassword: hashedPassword,
		Firstname:      "John",
		Lastname:       "Doe",
		DisplayName:    "John Doe",
		BirthDate:      birthDate,
		Gender:         entity.Male,
	}

	return user
}
