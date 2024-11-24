package app

import (
	"tg-service-v2/internal/api/delivery/http/status"
	"tg-service-v2/internal/api/delivery/telegram/cars"
	"tg-service-v2/internal/api/delivery/telegram/start"
	"tg-service-v2/internal/api/delivery/telegram/users"
	"tg-service-v2/internal/api/delivery/telegram/watcher"
)

func (a *App) initHandlers() {
	a.statusHandler = status.NewHandler(a.logger, a.serviceName)

	a.userHandler = users.NewHandler(a.userService, a.redisService, a.userMapsService, a.chain, a.logger)
	a.carHandler = cars.NewHandler(a.logger, a.tgBot, a.config.Extra, a.carService, a.redisService, a.userMapsService)
	a.watcherHandler = watcher.NewHandler(a.userService, a.redisService, a.userMapsService, a.carService, a.chain, a.logger)
	a.startHandler = start.NewHandler(a.logger)

	a.logger.Debug("handlers created")
}
