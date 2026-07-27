package app

import (
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/adapters"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

type RepositoryConstructor struct {
	Builder func(app *App) repositories.BaseRepository
}

type ServiceConstructor func(app *App) interface{}

type ControllerConstructor func(app *App)

type TypedRepositories struct {
	UserRepository  *repositories.UserRepository
	ImageRepository *repositories.ImageRepository
	FileStorage     *repositories.LocalStorage
}

type TypedServices struct {
	Healthcheck  *services.HealthCheckService
	UserService  *services.UserService
	ImageService *services.ImageService
}

type TypedAdapters struct {
	PasswordHasher *adapters.Argon2Hasher
	TokenProvider  *adapters.JWTProvider
}
