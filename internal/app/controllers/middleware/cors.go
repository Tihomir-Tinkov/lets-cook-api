package middleware

import (
	"net/http"
	"strings"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/config"
)

func CorsMiddleware(next http.Handler, corsConfig config.CorsConfig) http.HandlerFunc {

	allowedReqHeaders := []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"X-Client-Id",
		"X-Entity-Limit",
		"X-User-Id",
	}

	allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"

	return func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if corsConfig.AllowOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			allowedOrigins := strings.Split(corsConfig.AllowOrigins, ",")
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedReqHeaders, ", "))

		if r.Method == http.MethodOptions {
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}
