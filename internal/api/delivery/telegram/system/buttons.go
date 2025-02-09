package system

import (
	"context"
	"github.com/andReyM228/lib/log"
	"gopkg.in/telebot.v3"
	"strconv"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/domain/menu"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	log             log.Logger
	tgBot           *telebot.Bot
	redisService    services.RedisService
	userMapsService services.UserMapsService
}

func NewHandler(
	log log.Logger,
	tgBot *telebot.Bot,
	redisService services.RedisService,
	userMapsService services.UserMapsService) Handler {

	menu.InitMainMenu()
	menu.InitCarsMenu()
	return Handler{
		log:             log,
		tgBot:           tgBot,
		redisService:    redisService,
		userMapsService: userMapsService,
	}
}

func (h Handler) BackButton() (*telebot.Btn, func(ctx telebot.Context) error) {
	return &menu.BackButton, func(ctx telebot.Context) error {
		_, err := h.redisService.GetToken(ctx.Chat().ID)
		if err != nil {
			h.log.Errorf("get user token error: %v", err)

			if err := ctx.Send("you are not authorized"); err != nil {
				h.log.Errorf("error send error message: %v", err)
				return err
			}

			return nil
		}

		if h.userMapsService.Fetch(context.Background(), domain.GenKey(domain.BuyersUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10))) {
			err := h.userMapsService.Delete(context.Background(), domain.GenKey(domain.BuyersUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)))
			if err != nil {
				h.log.Errorf("delete buy-user error: %v", err)
				if err := ctx.Send("you are not authorized"); err != nil {
					h.log.Errorf("error send error message: %v", err)
					return err
				}
				return err
			}

			if err := ctx.Send(" ℹ️ Return to cars-menu ℹ️ ", menu.CarsMenu); err != nil {
				return err
			}

			return nil
		}

		if err := ctx.Send(" ℹ️ Return to start-menu ℹ️ ", menu.MainMenu); err != nil {
			return err
		}

		return nil
	}
}
