package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/ports"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization format")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
)

type AuthMiddleware struct {
	tokens ports.TokenProvider
}

func NewAuthMiddleware(
	tokens ports.TokenProvider,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokens: tokens,
	}
}

func extractBearerToken(header string) (string, bool) {
	parts := strings.Fields(header)

	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	return parts[1], true
}

func (m *AuthMiddleware) RequireAuth(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			responses.JSONError(
				w,
				r,
				ErrMissingAuthHeader,
				http.StatusUnauthorized,
			)
			return
		}

		token, ok := extractBearerToken(authHeader)

		if !ok {
			responses.JSONError(
				w,
				r,
				ErrInvalidAuthFormat,
				http.StatusUnauthorized,
			)
			return
		}

		authContext, err := m.tokens.Validate(token)

		if err != nil {
			responses.JSONError(
				w,
				r,
				ErrInvalidToken,
				http.StatusUnauthorized,
			)
			return
		}

		ctx := SetAuthContext(
			r.Context(),
			authContext,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func (m *AuthMiddleware) RequireRole(
	role models.UserRole,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			auth, ok := GetAuthContext(
				r.Context(),
			)

			if !ok {
				responses.JSONError(
					w,
					r,
					ErrUnauthorized,
					http.StatusUnauthorized,
				)
				return
			}

			if auth.Role != role {
				responses.JSONError(
					w,
					r,
					ErrForbidden,
					http.StatusForbidden,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
