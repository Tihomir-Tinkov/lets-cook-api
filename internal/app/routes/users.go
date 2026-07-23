package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterUserRoutes(
	r *Router,
	ctl *controllers.UserController,
) {
	userRoutes := []Route{
		{
			Path: "/users/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet: ctl.GetUser,
			},
		},
	}

	r.routes = append(r.routes, userRoutes...)
}
