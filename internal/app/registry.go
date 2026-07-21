// app/registry.go
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
	//Zones              *repositories.ZoneRepository
	//GroupZones         *repositories.GroupZoneRepository
	//Objects            *repositories.ObjectRepository
	//Roles              *repositories.RoleRepository
	//Notifications      *repositories.NotificationRepository
	//Permissions        *repositories.PermissionRepository
	//Sensors            *repositories.SensorRepository
	//GroupAlarms        *repositories.GroupAlarmRepository
	//ObjectDevice       *repositories.ObjectDeviceRepository
	//ObjectMobilisights *repositories.ObjectMobilisightsRepository
	//ObjectJimi         *repositories.JimiRepository
	//JimiAlarmLogs      *repositories.JimiAlarmLogRepository
	//JimiInstructLogs   *repositories.JimiInstructLogRepository
}

type TypedServices struct {
	Healthcheck *services.HealthCheckService
}
