package router

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/handler"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/middleware" // Import middleware
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	// === Inisialisasi Repository dan Handler ===
	// User/Auth
	userRepository := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepository)
	// Product
	productRepository := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepository)

	orderRepository := repository.NewOrderRepository(db)
	orderHandler := handler.NewOrderHandler(orderRepository, productRepository)

	// Grup rute API v1
	api := app.Group("/api/v1")

	// Rute untuk Auth (Public)
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// Rute untuk Produk
	// Rute GET bisa diakses publik
	api.Get("/products", productHandler.GetAllProducts)
	api.Get("/products/:id", productHandler.GetProductByID)

	// Rute yang membutuhkan otentikasi
	protected := api.Group("/", middleware.AuthRequired())

	// Hanya admin yang bisa membuat, mengubah, dan menghapus produk
	protected.Post("/products", middleware.AdminOnly(), productHandler.CreateProduct)
	protected.Put("/products/:id", middleware.AdminOnly(), productHandler.UpdateProduct)
	protected.Delete("/products/:id", middleware.AdminOnly(), productHandler.DeleteProduct)

	// Rute Pesanan untuk semua user yang terotentikasi
	protected.Post("/orders", orderHandler.CreateOrder)
	protected.Get("/orders", orderHandler.GetUserOrders)
	protected.Get("/orders/:id", orderHandler.GetOrderByID)
}