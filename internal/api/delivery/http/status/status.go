package status

import (
	"github.com/andReyM228/lib/log"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Handler struct {
	log         log.Logger
	serviceName string
	date        time.Time
}

func NewHandler(log log.Logger, serviceName string) Handler {
	return Handler{
		log:         log,
		serviceName: serviceName,
		date:        time.Now().UTC(),
	}
}

func (h Handler) CheckStatus(ctx *fiber.Ctx) error {
	return ctx.JSON(response{
		Name:      h.serviceName,
		BuildDate: h.date,
	})
}
