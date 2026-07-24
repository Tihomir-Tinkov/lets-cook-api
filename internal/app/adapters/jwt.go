package adapters

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: check correctness

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTProvider struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTProvider(secret string, expiration time.Duration) *JWTProvider {
	return &JWTProvider{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (j *JWTProvider) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(j.expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}

func (j *JWTProvider) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			// Prevent algorithm switching attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return j.secret, nil
		},
	)

	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	userIDString, ok := claims["user_id"].(string)

	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDString)

	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}
