package cars

import (
	"github.com/andReyM228/lib/log"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	log          log.Logger
	carService   services.CarService
	redisService services.RedisService
}

func NewHandler(log log.Logger, carService services.CarService, redisService services.RedisService) Handler {
	initMainMenu()

	return Handler{
		log:          log,
		carService:   carService,
		redisService: redisService,
	}
}

func (h Handler) GetUserCars(token string) (domain.Cars, error) {
	cars, err := h.carService.GetUserCars(token)
	if err != nil {
		return domain.Cars{}, err
	}

	return cars, nil
}

func (h Handler) BuyCar(token, txHash string, carID int64) error {
	err := h.carService.BuyCar(token, txHash, carID)
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) SellCar(chatID, carID int64, token string) error {
	err := h.carService.SellCar(chatID, carID, token)
	if err != nil {
		return err
	}

	return nil
}
