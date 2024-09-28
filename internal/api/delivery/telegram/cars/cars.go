package cars

import (
	"bytes"
	"encoding/base64"
	"github.com/andReyM228/lib/log"
	"gopkg.in/telebot.v3"
	"strconv"
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

func (h Handler) GetCar(ctx telebot.Context) error {
	chatID := ctx.Sender().ID

	token, err := h.redisService.GetToken(chatID)
	if err != nil {
		h.log.Errorf("/get-car error: ", err)

		err := ctx.Send("error while get car")
		if err != nil {
			h.log.Error("failed ctx.Send")
		}

		return nil
	}

	carIDString := ctx.Message().Payload

	carID, err := strconv.Atoi(carIDString)
	if err != nil {
		h.log.Errorf("/get-car error (convert carID to int64): ", err)

		err := ctx.Send("error while get car")
		if err != nil {
			h.log.Error("failed ctx.Send")
		}

		return nil
	}

	car, err := h.carService.GetCar(int64(carID), token)
	if err != nil {
		return err
	}

	resp := showCar(car)

	carImageBytes, err := base64.StdEncoding.DecodeString(car.Image)
	if err != nil {
		h.log.Errorf("/get-car error (decode base64 image): ", err)

		err := ctx.Send("error while get car")
		if err != nil {
			h.log.Error("failed ctx.Send")
		}

		return nil
	}

	reader := bytes.NewReader(carImageBytes)

	carImage := &telebot.Photo{
		File:    telebot.FromReader(reader),
		Caption: resp,
	}

	err = ctx.SendAlbum(telebot.Album{carImage})
	if err != nil {
		h.log.Errorf("/get-car error (send resp): ", err)

		err := ctx.Send("error while get car")
		if err != nil {
			h.log.Error("failed ctx.Send")
		}

		return nil
	}

	return nil
}
