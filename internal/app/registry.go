package app

import (
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/adapters"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/repositories"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/services"
)

type RepositoryConstructor struct {
	Builder func(app *App) repositories.BaseRepository
}

type ServiceConstructor func(app *App) interface{}

type ControllerConstructor func(app *App)

type TypedRepositories struct {
	UserRepository    *repositories.UserRepository
	RecipeRepository  *repositories.RecipeRepository
	CommentRepository *repositories.CommentRepository
	ImageRepository   *repositories.ImageRepository
	FileStorage       *repositories.LocalStorage
}

type TypedServices struct {
	Healthcheck    *services.HealthCheckService
	UserService    *services.UserService
	RecipeService  *services.RecipeService
	CommentService *services.CommentService
	ImageService   *services.ImageService
}

type TypedAdapters struct {
	PasswordHasher *adapters.Argon2Hasher
	TokenProvider  *adapters.JWTProvider
}
