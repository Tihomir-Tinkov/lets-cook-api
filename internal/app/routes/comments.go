package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterCommentRoutes(r *Router, ctl *controllers.CommentController) {
	commentRoutes := []Route{
		{
			Path: "/recipes/{recipeId}/comments",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet:  ctl.GetByRecipeID,
				http.MethodPost: ctl.Create,
			},
		},
		{
			Path: "/recipes/{recipeId}/comments/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet:    ctl.GetByID,
				http.MethodPut:    ctl.Update,
				http.MethodDelete: ctl.Delete,
			},
		},
	}

	r.routes = append(r.routes, commentRoutes...)
}
