package app

import (
	"tg-service-v2/internal/api/delivery/http/status"
	"tg-service-v2/internal/api/delivery/telegram/users"
)

func (a *App) initHandlers() {
	a.statusHandler = status.NewHandler(a.logger, a.serviceName)

	a.userHandler = users.NewHandler(a.userService, a.chain)

	a.logger.Debug("handlers created")
}
