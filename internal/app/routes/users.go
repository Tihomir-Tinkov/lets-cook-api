package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterUserRoutes(r *Router, ctl *controllers.UserController) {
	userRoutes := []Route{
		{
			Path: "/auth/register",
			Methods: map[string]http.HandlerFunc{
				http.MethodPost: ctl.Register,
			},
		},
		{
			Path: "/auth/login",
			Methods: map[string]http.HandlerFunc{
				http.MethodPost: ctl.Login,
			},
		},
		{
			Path: "/auth/logout",
			Methods: map[string]http.HandlerFunc{
				http.MethodPost: ctl.Logout,
			},
		},
		{
			Path: "/users/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet: ctl.GetByID,
			},
		},
		{
			Path: "/users/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodPut: ctl.Update,
			},
		},
		{
			Path: "/users/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodDelete: ctl.Delete,
			},
		},
	}

	r.routes = append(r.routes, userRoutes...)
}
