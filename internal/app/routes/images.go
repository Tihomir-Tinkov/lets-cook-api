package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterImagesRoutes(r *Router, ctl *controllers.ImageController) {
	imagesRoutes := []Route{
		{
			Path: "/images",
			Methods: map[string]http.HandlerFunc{
				http.MethodPost: ctl.Upload,
			},
		},
		{
			Path: "/images/{id}",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet:    ctl.Download,
				http.MethodDelete: ctl.Delete,
			},
		},
	}

	r.routes = append(r.routes, imagesRoutes...)
}
