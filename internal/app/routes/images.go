package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers/middleware"
)

func RegisterImagesRoutes(r *Router, ctl *controllers.ImageController, authMiddleware *middleware.AuthMiddleware) {
	imagesRoutes := []Route{
		{
			Path: "/images",
			Methods: map[string]Handler{
				http.MethodPost: {
					Func: ctl.Upload,
					Middlewares: []Middleware{
						authMiddleware.RequireAuth,
					},
				},
			},
		},
		{
			Path: "/images/{id}",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.Download,
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

	r.routes = append(r.routes, imagesRoutes...)
}
