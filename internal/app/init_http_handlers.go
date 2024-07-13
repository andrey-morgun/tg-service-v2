package app

import "tg-service-v2/internal/api/delivery/http/status"

func (a *App) initHandlers() {
	a.statusHandler = status.NewHandler(a.logger, a.serviceName)

	a.logger.Debug("handlers created")
}
