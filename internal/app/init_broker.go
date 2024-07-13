package app

import (
	"github.com/andReyM228/lib/rabbit"
)

func (a *App) initBroker() {
	var err error
	a.rabbit, err = rabbit.NewRabbitMQ(a.config.Rabbit.RabbitUrl, a.logger)
	if err != nil {
		a.logger.Fatal(err.Error())
	}

}
