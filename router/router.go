package router

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/handler"             // Path ke handler
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository" // Path ke repository
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

// SetupRoutes mendefinisikan dan mengatur semua rute untuk aplikasi
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Middleware
	app.Use(logger.New())

	// === Inisialisasi Repository dan Handler ===
	productRepository := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepository)

	// Grup rute untuk API v1
	api := app.Group("/api/v1")

	// Rute untuk Produk
	api.Get("/products", productHandler.GetAllProducts)
	api.Get("/products/:id", productHandler.GetProductByID)
	api.Post("/products", productHandler.CreateProduct)
	api.Put("/products/:id", productHandler.UpdateProduct)
	api.Delete("/products/:id", productHandler.DeleteProduct)
}