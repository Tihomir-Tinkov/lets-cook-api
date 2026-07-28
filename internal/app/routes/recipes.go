package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/middleware"
)

func RegisterRecipeRoutes(r *Router, ctl *controllers.RecipeController, authMiddleware *middleware.AuthMiddleware) {
	recipeRoutes := []Route{
		{
			Path: "/recipes/",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.List,
				},
				http.MethodPost: {
					Func: ctl.Create,
					Middlewares: []Middleware{
						authMiddleware.RequireAuth,
					},
				},
			},
		},
		{
			Path: "/recipes/{id}",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.GetByID,
				},
				http.MethodPut: {
					Func: ctl.Update,
					Middlewares: []Middleware{
						authMiddleware.RequireAuth,
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

	r.routes = append(r.routes, recipeRoutes...)
}
