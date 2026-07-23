package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/middleware"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/rs/zerolog/log"
)

type Router struct {
	mux                     *http.ServeMux
	notFoundHandler         http.HandlerFunc
	methodNotAllowedHandler http.HandlerFunc
	prefix                  string
	routes                  []Route
	loggerMiddleware        bool
}

type Route struct {
	Methods map[string]http.HandlerFunc
	Path    string
}

func NewRouter(opts ...RouterOptFunc) *Router {
	r := &Router{
		mux: http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

func (r *Router) RegisterRoutes() {
	r.registerRoutes()
}

func (r *Router) registerRoutes() {

	r.notFoundHandler = func(w http.ResponseWriter, req *http.Request) {
		responses.JSONError(w, req, errors.New("not found"), http.StatusNotFound)
	}

	r.methodNotAllowedHandler = func(w http.ResponseWriter, req *http.Request) {
		responses.JSONError(w, req, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}

	for _, rt := range r.routes {
		route := rt
		path := route.Path
		if r.prefix != "" {
			path = fmt.Sprintf("/%s%s", r.prefix, path)
		}

		for method, handler := range route.Methods {
			pattern := fmt.Sprintf("%s %s", method, path)

			r.mux.HandleFunc(pattern, r.recoverer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if r.loggerMiddleware {
					middleware.LoggerMiddleware(handler)(w, req)
					return
				}

				handler(w, req)
			})))
		}
	}
}

func (r *Router) recoverer(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().Interface("panic", rec).Msg("handler panic recovered")
				responses.JSONError(w, req, errors.New("internal server error"), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, req)
	}
}

func (r *Router) ListRoutes() {
	for _, route := range r.routes {
		for method := range route.Methods {
			fields := map[string]interface{}{
				"method": method,
				"route":  fmt.Sprintf("/%s%s", r.prefix, route.Path),
			}
			log.Debug().Fields(fields).Msg("")
		}
	}
}
