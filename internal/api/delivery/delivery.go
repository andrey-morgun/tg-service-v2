package delivery

import "github.com/gofiber/fiber/v2"

type (
	StatusHandler interface {
		CheckStatus(ctx *fiber.Ctx) error
	}

	BrokerStatusHandler interface {
		CheckStatus(request []byte) error
	}
)
