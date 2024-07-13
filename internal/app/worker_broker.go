package app

import (
	"context"
)

func serveBroker(ctx context.Context, a *App) {
	a.initBrokerHandlers()

	// graceful shutdown listener.
	go func() {
		<-ctx.Done()

		if err := a.rabbit.Close(); err != nil {
			a.logger.Infof("🔵 rabbit: broker close: %v", err)
		}
	}()

	err := a.rabbit.ConsumeAll(ctx, a.registerBrokerTopics())
	if err != nil {
		a.logger.Fatal(err.Error())
	}
}
