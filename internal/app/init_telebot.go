package app

import (
	"gopkg.in/telebot.v3"
	"time"
)

func (a *App) initTelebot() {
	var err error
	a.tgBot, err = telebot.NewBot(telebot.Settings{Token: a.config.TgBot.Token, Poller: &telebot.LongPoller{Timeout: 10 * time.Second}, ParseMode: telebot.ModeMarkdown})
	if err != nil {
		a.logger.Fatal(err.Error())
		return
	}
}
