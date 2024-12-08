package app

import (
	"context"
	"gopkg.in/telebot.v3"
)

func (a *App) registerTelebotCommands(ctx context.Context) {
	a.tgBot.Handle("/start", a.startHandler.Start)

	a.tgBot.Handle("/registration", a.userHandler.Registration)
	a.tgBot.Handle("/login", a.userHandler.Login)
	a.tgBot.Handle("/getcar", a.carHandler.GetCar)
	a.tgBot.Handle(telebot.OnText, a.watcherHandler.MsgWatcher)

	a.tgBot.Handle(a.carHandler.GetCarsButton())
	a.tgBot.Handle(a.carHandler.BuyCarButton())
	a.tgBot.Handle(a.carHandler.GetCarsMenu())
	a.tgBot.Handle(a.carHandler.UserCarsButton())
}
