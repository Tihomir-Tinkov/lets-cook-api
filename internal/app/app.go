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

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
	gateways "github.com/Tihomir-Tinkov/cooking-site-project/internal/app/gateway"
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
	DB            *pgxpool.Pool
	Router        *routes.Router
	Server        *http.Server
	TypedRepos    *TypedRepositories
	TypedServices *TypedServices
	Cache         *cache.TokenCache
	Config        config.Config
	closeOnce     sync.Once
	closeErr      error
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

	app.bootstrap()

	return app, nil
}

// controller constructors remain for route wiring

var controllerConstructors = map[string]ControllerConstructor{
	"healthcheck": func(app *App) {
		ctl := controllers.NewServiceController(app.TypedServices.Healthcheck)
		routes.RegisterServiceRoutes(app.Router, ctl)
	},
}

func (a *App) bootstrap() {
	// Instantiate typed repositories directly
	//a.TypedRepos = &TypedRepositories{
	//	Zones:              repositories.NewZonesRepository(a.DB.GPS),
	//	GroupZones:         repositories.NewGroupZoneRepository(a.DB.GPS),
	//	Objects:            repositories.NewObjectRepository(a.DB.GPS).(*repositories.ObjectRepository),
	//	Roles:              repositories.NewRoleRepository(a.DB.GPS),
	//	Notifications:      repositories.NewNotificationRepository(a.DB.GPS),
	//	Permissions:        repositories.NewPermissionRepository(a.DB.GPS),
	//	Sensors:            repositories.NewRepository(api.New(a.Config.API.Reporting)),
	//	GroupAlarms:        repositories.NewGroupAlarmRepository(api.New(a.Config.API.Alarms)),
	//	ObjectDevice:       repositories.NewObjectDeviceRepository(a.DB.GPS),
	//	ObjectMobilisights: repositories.NewObjectMobilisightsRepository(a.DB.GPS),
	//	ObjectJimi:         repositories.NewJimiRepository(a.DB.GPS, a.Redis.Jimi),
	//	JimiAlarmLogs:      repositories.NewJimiAlarmLogRepository(a.DB.GPS),
	//	JimiInstructLogs:   repositories.NewJimiInstructLogRepository(a.DB.GPS),
	//}

	healthProbe := gateways.NewHealth(a.DB)

	a.TypedServices = &TypedServices{
		Healthcheck: services.NewHealthCheckService(healthProbe),
	}

	for _, constructor := range controllerConstructors {
		constructor(a)
	}

	a.Router.RegisterRoutes()

	a.Router.ListRoutes()
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
