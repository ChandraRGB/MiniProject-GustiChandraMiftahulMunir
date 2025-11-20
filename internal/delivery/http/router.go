package http

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes registers all HTTP routes for the application.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// TODO: register other feature routes (auth, user, toko, alamat, kategori, produk, trx)
}
