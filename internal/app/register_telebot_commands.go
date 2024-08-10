package app

import (
	"context"
	"gopkg.in/telebot.v3"
)

func (a *App) registerTelebotCommands(ctx context.Context) {
	a.tgBot.Handle("/start", func(c telebot.Context) error {
		return nil
	})

	a.tgBot.Handle("/registration", a.userHandler.Registration)
	a.tgBot.Handle("/login", a.userHandler.Login)
	a.tgBot.Handle(telebot.OnText, a.userHandler.MsgWatcher)

	a.tgBot.Handle(a.carHandler.GetCarsButton())
}
