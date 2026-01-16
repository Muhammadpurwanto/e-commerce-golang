package main

import (
	"log"

	"github.com/Muhammadpurwanto/e-commerce-golang/config"
	"github.com/Muhammadpurwanto/e-commerce-golang/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Memuat Konfigurasi
	port := config.LoadConfig()

	// 2. Inisialisasi Fiber App
	app := fiber.New()

	// 3. Setup Rute
	router.SetupRoutes(app)

	// 4. Menjalankan Server
	log.Printf("Server is starting on port %s\n", port)
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}