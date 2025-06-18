package routes

import (
	"go-url-shortener/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/:url", handler.ResolveUrl)
	app.Post("/api/v1", handler.ShortenURL)
}
