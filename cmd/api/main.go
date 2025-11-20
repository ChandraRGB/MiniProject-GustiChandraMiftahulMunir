package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/config"
	httpDelivery "github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/delivery/http"
)

func main() {
	// init database connection
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// init Fiber app
	app := fiber.New()

	// register routes
	httpDelivery.RegisterRoutes(app, db)

	// start server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
