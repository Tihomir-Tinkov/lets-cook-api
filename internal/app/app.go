package app

import (
	"context"
	"errors"
	"net/http"
	"net/netip"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/adapters"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/middleware"
	gateways "github.com/Tihomir-Tinkov/cooking-site-project/internal/app/gateway"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/routes"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/cache"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/config"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/conn/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type App struct {
	DB             *pgxpool.Pool
	Router         *routes.Router
	AuthMiddleware *middleware.AuthMiddleware
	Server         *http.Server
	TypedRepos     *TypedRepositories
	TypedAdapters  *TypedAdapters
	TypedServices  *TypedServices
	Cache          *cache.TokenCache
	Config         config.Config
	closeOnce      sync.Once
	closeErr       error
}

func NewApp(cfg config.Config) (*App, error) {
	dbConn, err := db.Connect(cfg.DB)
	if err != nil {
		return nil, err
	}

	// Build router with optional logger middleware
	var opts []routes.RouterOptFunc
	if cfg.Routes.LoggerMiddleware {
		opts = append(opts, routes.WithLoggingMiddleware())
	}

	opts = append(opts, routes.WithPrefix(cfg.Routes.Prefix))

	app := &App{
		Config: cfg,
		Router: routes.NewRouter(opts...),
	}
	app.DB = dbConn

	addr := netip.AddrPortFrom(netip.IPv4Unspecified(), cfg.Server.InternalPort)

	app.Server = &http.Server{
		Addr:              addr.String(),
		Handler:           app.Router.Mux(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	app.Cache = cache.NewTokenCache()
	app.Cache.Cleanup(1 * time.Hour) // lower the time according to the fastest expiring token

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}

// controller constructors remain for route wiring

var controllerConstructors = map[string]ControllerConstructor{
	"healthcheck": func(app *App) {
		ctl := controllers.NewServiceController(app.TypedServices.Healthcheck)
		routes.RegisterServiceRoutes(app.Router, ctl)
	},

	"user": func(app *App) {
		ctl := controllers.NewUserController(app.TypedServices.UserService)
		routes.RegisterUserRoutes(app.Router, ctl, app.AuthMiddleware)
	},

	"recipe": func(app *App) {
		ctl := controllers.NewRecipeController(app.TypedServices.RecipeService)
		routes.RegisterRecipeRoutes(app.Router, ctl, app.AuthMiddleware)
	},

	"comment": func(app *App) {
		ctl := controllers.NewCommentController(app.TypedServices.CommentService)
		routes.RegisterCommentRoutes(app.Router, ctl, app.AuthMiddleware)
	},

	"images": func(app *App) {
		ctl := controllers.NewImageController(app.TypedServices.ImageService)
		routes.RegisterImagesRoutes(app.Router, ctl, app.AuthMiddleware)
	},
}

func (a *App) bootstrap() error {
	// Instantiate typed repositories directly
	a.TypedRepos = &TypedRepositories{
		UserRepository:    repositories.NewUserRepository(a.DB),
		RecipeRepository:  repositories.NewRecipeRepository(a.DB),
		CommentRepository: repositories.NewCommentRepository(a.DB),
		ImageRepository:   repositories.NewImageRepository(a.DB),
		FileStorage:       repositories.NewLocalStorage(a.Config.StorePath),
	}

	healthProbe := gateways.NewHealth(a.DB)

	hasher, err := adapters.NewArgon2Hasher(adapters.Argon2Config(a.Config.Argon2))

	if err != nil {
		return err
	}

	jwtprovider, err := adapters.NewJWTProvider(a.Config.JWT.Secret, a.Config.JWT.Issuer, a.Config.JWT.Expiration)

	if err != nil {
		return err
	}

	a.TypedAdapters = &TypedAdapters{
		PasswordHasher: hasher,
		TokenProvider:  jwtprovider,
	}

	a.AuthMiddleware = middleware.NewAuthMiddleware(a.TypedAdapters.TokenProvider)

	a.TypedServices = &TypedServices{
		Healthcheck:    services.NewHealthCheckService(healthProbe),
		UserService:    services.NewUserService(a.TypedRepos.UserRepository, a.TypedAdapters.PasswordHasher, a.TypedAdapters.TokenProvider),
		RecipeService:  services.NewRecipeService(a.TypedRepos.RecipeRepository),
		CommentService: services.NewCommentService(a.TypedRepos.CommentRepository),
		ImageService:   services.NewImageService(a.TypedRepos.ImageRepository, a.TypedRepos.FileStorage),
	}

	for _, constructor := range controllerConstructors {
		constructor(a)
	}

	a.Router.RegisterRoutes()

	a.Router.ListRoutes()

	return nil
}

func (a *App) ServeAndListen() error {
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-shutdownCtx.Done()
	log.Info().Msg("shutdown signal received, shutting down...")

	return nil
}

func (a *App) Close() error {
	a.closeOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := a.Server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("graceful shutdown failed")
			if closeErr := a.Server.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
				a.closeErr = closeErr
				return
			}
		}

		var group errgroup.Group
		group.Go(func() error {
			a.DB.Close()
			return nil
		})

		a.closeErr = group.Wait()
	})

	return a.closeErr
}
