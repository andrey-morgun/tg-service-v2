package cars

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/andReyM228/lib/log"
	"gopkg.in/telebot.v3"
	"strconv"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/domain/menu"
	"tg-service-v2/internal/api/services"
	"tg-service-v2/internal/config"
)

type Handler struct {
	log             log.Logger
	tgBot           *telebot.Bot
	config          config.Extra
	carService      services.CarService
	redisService    services.RedisService
	userMapsService services.UserMapsService
}

func NewHandler(
	log log.Logger,
	tgBot *telebot.Bot,
	config config.Extra,
	carService services.CarService,
	redisService services.RedisService,
	userMapsService services.UserMapsService) Handler {

	menu.InitMainMenu()
	menu.InitCarsMenu()
	menu.InitTransferMenu()
	return Handler{
		log:             log,
		tgBot:           tgBot,
		config:          config,
		carService:      carService,
		redisService:    redisService,
		userMapsService: userMapsService,
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

	car, err := h.carService.GetCar(context.Background(), int64(carID), token)
	if err != nil {
		return err
	}

	resp := showCar(car)

	reader := bytes.NewReader(car.ImageBytes)

	carImage := &telebot.Photo{
		File:    telebot.FromReader(reader),
		Caption: resp,
	}

	carDataJson, err := json.Marshal(domain.CarIDAndPrice{
		ID:    car.ID,
		Price: car.Price,
	})
	if err != nil {
		return err
	}

	buyCarBtn := telebot.Btn{
		Unique: "buy_car",
		Text:   "Buy Car",
		Data:   string(carDataJson),
	}

	inlineKeyboard := &telebot.ReplyMarkup{}
	inlineKeyboard.Inline(
		inlineKeyboard.Row(buyCarBtn),
	)

	err = ctx.Send(carImage, inlineKeyboard)
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
