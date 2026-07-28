package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/middleware"
)

func RegisterCommentRoutes(r *Router, ctl *controllers.CommentController, authMiddleware *middleware.AuthMiddleware) {
	commentRoutes := []Route{
		{
			Path: "/recipes/{recipeId}/comments",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.GetByRecipeID,
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
			Path: "/recipes/{recipeId}/comments/{id}",
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

	r.routes = append(r.routes, commentRoutes...)
}
