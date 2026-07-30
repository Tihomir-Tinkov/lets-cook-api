package middleware

import (
	"context"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
)

type contextKey string

const authContextKey contextKey = "auth"

func SetAuthContext(
	ctx context.Context,
	auth *models.AuthContext,
) context.Context {
	return context.WithValue(
		ctx,
		authContextKey,
		auth,
	)
}

func GetAuthContext(
	ctx context.Context,
) (*models.AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey).(*models.AuthContext)

	return auth, ok
}
