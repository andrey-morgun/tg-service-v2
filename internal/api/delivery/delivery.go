package delivery

import (
	"github.com/gofiber/fiber/v2"
	"gopkg.in/telebot.v3"
)

type (
	StartHandler interface {
		Start(ctx telebot.Context) error
	}

	StatusHandler interface {
		CheckStatus(ctx *fiber.Ctx) error
	}

	BrokerStatusHandler interface {
		CheckStatus(request []byte) error
	}

	UserHandler interface {
		Registration(ctx telebot.Context) error
		Login(ctx telebot.Context) (err error)
	}

	CarHandler interface {
		GetCarsButton() (*telebot.Btn, func(ctx telebot.Context) error)
		BuyCarButton() (*telebot.Btn, func(ctx telebot.Context) error)
		GetCar(ctx telebot.Context) error
		GetCarsMenu() (*telebot.Btn, func(ctx telebot.Context) error)
		UserCarsButton() (*telebot.Btn, func(ctx telebot.Context) error)
	}

	SystemButtons interface {
		BackButton() (*telebot.Btn, func(ctx telebot.Context) error)
	}

	Watcher interface {
		MsgWatcher(ctx telebot.Context) (err error)
	}
)
