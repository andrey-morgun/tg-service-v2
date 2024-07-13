package app

import "tg-service-v2/internal/api/delivery/broker/status"

func (a *App) initBrokerHandlers() {
	a.statusBrokerHandler = status.NewHandler(a.logger, a.serviceName, a.rabbit)

	a.logger.Debug("handlers created")
}
