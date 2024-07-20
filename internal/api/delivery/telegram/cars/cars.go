package cars

import (
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	carService services.CarService
}

func NewHandler(carService services.CarService) Handler {
	return Handler{
		carService: carService,
	}
}

//func (h Handler) GetCar(ctx telebot.Context) error {
//	msg := ctx.Message().Text
//
//	car, err := h.carService.GetCar(id, token)
//	if err != nil {
//		return err
//	}
//
//	return ctx.Send(car)
//}

func (h Handler) GetAllCars(token, label string) (domain.Cars, error) {
	cars, err := h.carService.GetCars(token, label)
	if err != nil {
		return domain.Cars{}, err
	}

	return cars, nil
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
