package middleware

import (
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
)

type AuthMiddleware struct {
	tokens ports.TokenProvider
}

func NewAuthMiddleware(tokens ports.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		tokens: tokens,
	}
}
