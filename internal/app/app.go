package app

import (
	"context"
	"embed"
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/lib/rabbit"
	"github.com/andReyM228/lib/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"tg-service-v2/internal/api/delivery"
	"tg-service-v2/internal/api/repository"
	"tg-service-v2/internal/api/repository/cars"
	redisRepo "tg-service-v2/internal/api/repository/redis"
	"tg-service-v2/internal/api/repository/users"
	"tg-service-v2/internal/api/service"
	"tg-service-v2/internal/api/service/car"
	redisService "tg-service-v2/internal/api/service/redis"
	"tg-service-v2/internal/api/service/user"
	"tg-service-v2/internal/config"
)

type (
	App struct {
		config config.Config
		//service     service.Service
		carRepo             repository.CarRepo
		userRepo            repository.UserRepo
		redisRepo           repository.RedisRepo
		carService          service.CarService
		userService         service.UserService
		redisService        service.RedisService
		statusHandler       delivery.StatusHandler
		statusBrokerHandler delivery.BrokerStatusHandler
		rabbit              rabbit.Rabbit
		logger              log.Logger
		validator           *validator.Validate
		serviceName         string
		router              *fiber.App
		redis               redis.Redis
		clientHTTP          *http.Client
	}
	worker func(ctx context.Context, a *App)
)

func New(name string) App {
	return App{
		serviceName: name,
	}
}

func (a *App) Run(ctx context.Context, fs embed.FS) {
	a.initLogger()
	a.initValidator()
	a.populateConfig()
	a.initBroker()
	a.initRepos()
	a.initServices()
	a.initHandlers()

	a.runWorkers(ctx)
}

func (a *App) initLogger() {
	a.logger = log.Init()
}

func (a *App) initRepos() {
	a.carRepo = cars.NewRepository(a.logger, a.clientHTTP, a.rabbit, a.config)
	a.userRepo = users.NewRepository(a.logger, a.clientHTTP, a.rabbit, a.validator)
	a.redisRepo = redisRepo.NewRepository(a.redis, a.logger)

	a.logger.Debug("repos created")
}

func (a *App) initServices() {
	a.carService = car.NewService(a.carRepo, a.logger)
	a.userService = user.NewService(a.userRepo, a.logger)
	a.redisService = redisService.NewService(a.redisRepo, a.logger)

	a.logger.Debug("services created")
}

func (a *App) initHTTPClient() {
	a.clientHTTP = http.DefaultClient
}

func (a *App) initRedis() {
	a.redis = redis.NewCache(a.config.Redis, a.logger)
}

func (a *App) initValidator() {
	a.validator = validator.New()
}

func (a *App) populateConfig() {
	cfg, err := config.ParseConfig()
	if err != nil {
		a.logger.Debugf("populateConfig: %s", err)
	}

	err = cfg.ValidateConfig(a.validator)
	if err != nil {
		a.logger.Debugf("populateConfig: %s", err)
	}

	a.config = cfg
}
