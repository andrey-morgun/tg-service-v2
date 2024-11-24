package start

import (
	"github.com/andReyM228/lib/log"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	log log.Logger
}

func NewHandler(log log.Logger) Handler {
	return Handler{
		log: log,
	}
}

func (h Handler) Start(ctx telebot.Context) error {
	if err := ctx.Send("Hi! write '/registration' or '/login' for start"); err != nil {
		return err
	}

	return nil
}
