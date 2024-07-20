package delivery

import (
	"github.com/gofiber/fiber/v2"
	"gopkg.in/telebot.v3"
)

type (
	StatusHandler interface {
		CheckStatus(ctx *fiber.Ctx) error
	}

	BrokerStatusHandler interface {
		CheckStatus(request []byte) error
	}

	UserHandler interface {
		Registration(ctx telebot.Context) error
	}
)
