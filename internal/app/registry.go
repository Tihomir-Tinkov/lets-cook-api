package app

import (
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

type RepositoryConstructor struct {
	Builder func(app *App) repositories.BaseRepository
}

type ServiceConstructor func(app *App) interface{}

type ControllerConstructor func(app *App)

type TypedRepositories struct {
	ImageRepository *repositories.ImageRepository
	FileStorage     *repositories.LocalStorage
}

type TypedServices struct {
	Healthcheck  *services.HealthCheckService
	ImageService *services.ImageService
}
