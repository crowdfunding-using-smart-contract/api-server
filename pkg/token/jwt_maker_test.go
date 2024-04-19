package token

import (
	"fund-o/api-server/pkg/random"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJwtTokenMaker(t *testing.T) {
	t.Run("Create and verify token", func(t *testing.T) {
		maker, err := NewJWTMaker(random.NewString(32))
		require.NoError(t, err)

		userID := uuid.NewString()
		duration := time.Minute

		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		token, payload, err := maker.CreateToken(userID, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		require.NotZero(t, payload.ID)
		require.Equal(t, userID, payload.UserID)
		require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
		require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	})
}

func TestExpiredJwtToken(t *testing.T) {
	t.Run("Expired token", func(t *testing.T) {
		maker, err := NewJWTMaker(random.NewString(32))
		require.NoError(t, err)

		token, payload, err := maker.CreateToken(uuid.NewString(), -time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrExpiredToken.Error())
		require.Nil(t, payload)
	})
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	t.Run("Invalid JWT token alg none", func(t *testing.T) {
		payload, err := NewPayload(uuid.NewString(), time.Minute)
		require.NoError(t, err)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
		token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		maker, err := NewJWTMaker(random.NewString(32))
		require.NoError(t, err)

		payload, err = maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrInvalidToken.Error())
		require.Nil(t, payload)
	})
}
