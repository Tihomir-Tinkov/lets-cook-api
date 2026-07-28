package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterRecipeRoutes(r *Router, ctl *controllers.RecipeController) {
	recipeRoutes := []Route{
		{
			Path: "/recipes/",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet:  ctl.List,
				http.MethodPost: ctl.Create,
			},
		},
		{
			Path: "/recipes/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet:    ctl.GetByID,
				http.MethodPut:    ctl.Update,
				http.MethodDelete: ctl.Delete,
			},
		},
	}

	r.routes = append(r.routes, recipeRoutes...)
}
