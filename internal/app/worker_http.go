package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func serveHttp(ctx context.Context, a *App) {
	a.router = fiber.New()

	a.registerHttpRoutes(a.router)

	// graceful shutdown listener.
	go func() {
		<-ctx.Done()

		if err := a.router.Shutdown(); err != nil {
			a.logger.Infof("🔵 http: server shutdown: %v", err)
		}
	}()

	a.logger.Debug("fiber api started")
	if err := a.router.Listen(fmt.Sprintf(":%d", a.config.HTTP.Port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatalf("🔴 failed to start server: %v", err)
		}
	}
}
