package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/middleware"
)

func RegisterUserRoutes(r *Router, ctl *controllers.UserController, authMiddleware *middleware.AuthMiddleware) {
	userRoutes := []Route{
		{
			Path: "/auth/register",
			Methods: map[string]Handler{
				http.MethodPost: {
					Func: ctl.Register,
				},
			},
		},
		{
			Path: "/auth/login",
			Methods: map[string]Handler{
				http.MethodPost: {
					Func: ctl.Login,
				},
			},
		},
		{
			Path: "/auth/logout",
			Methods: map[string]Handler{
				http.MethodPost: {
					Func: ctl.Logout,
				},
			},
		},
		{
			Path: "/users/{id}",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.GetByID,
				},
				http.MethodPut: {
					Func: ctl.Update,
					Middlewares: []Middleware{
						authMiddleware.RequireAuth,
						//authMiddleware.RequireRole(models.RoleAdmin),
					},
				},
				http.MethodDelete: {
					Func: ctl.Delete,
					Middlewares: []Middleware{
						authMiddleware.RequireAuth,
					},
				},
			},
		},
	}

	r.routes = append(r.routes, userRoutes...)
}
