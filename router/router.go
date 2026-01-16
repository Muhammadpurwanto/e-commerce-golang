package router

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes mendefinisikan dan mengatur semua rute untuk aplikasi
func SetupRoutes(app *fiber.App) {
	// Middleware untuk logging setiap request
	app.Use(logger.New())

	// Grup rute untuk API, misalnya /api
	api := app.Group("/api")

	// Grup rute untuk versi 1, misalnya /api/v1
	v1 := api.Group("/v1")

	// Menghubungkan rute GET "/" ke RootHandler
	v1.Get("/", handler.RootHandler)
}