package app

import "github.com/gofiber/fiber/v2"

func (a *App) registerHttpRoutes(app *fiber.App) {
	router := app.Group("/v1")

	router.Get("/status", a.statusHandler.CheckStatus)
}
