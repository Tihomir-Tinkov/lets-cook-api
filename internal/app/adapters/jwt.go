package adapters

import (
	"errors"
	"time"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	Role models.UserRole `json:"role"`

	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

func NewJWTProvider(secret string, issuer string, expiration time.Duration) (*JWTProvider, error) {
	if len(secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 bytes")
	}

	if issuer == "" {
		return nil, errors.New("jwt issuer required")
	}

	if expiration <= 0 {
		return nil, errors.New("jwt expiration must be positive")
	}
	return &JWTProvider{
		secret:     []byte(secret),
		issuer:     issuer,
		expiration: expiration,
	}, nil
}

func (j *JWTProvider) Generate(userID uuid.UUID, role models.UserRole) (string, error) {
	now := time.Now()

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiration)),
		},
	}

	// ECDSA, Ed25519, HMAC, RSA or RSAPSS
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *JWTProvider) Validate(tokenString string) (*models.AuthContext, error) {
	if len(tokenString) > 4096 {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrInvalidToken
			}
			return j.secret, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithLeeway(5*time.Second),
	)

	// log?
	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Subject == "" || claims.ExpiresAt == nil {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &models.AuthContext{
		UserID: userID,
		Role:   claims.Role,
	}, nil
}
