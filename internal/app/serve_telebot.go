package app

import (
	"context"
)

func serveTelebot(ctx context.Context, a *App) {
	a.registerTelebotCommands(ctx)

	// graceful shutdown listener.
	go func() {
		<-ctx.Done()

		if err := a.tgBot.Stop; err != nil {
			a.logger.Infof("🔵 tg-bot: server shutdown: %v", err)
		}
	}()

	a.tgBot.Start()
}
