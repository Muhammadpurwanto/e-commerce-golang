package main

import (
	"log"
	"os"

	"github.com/Muhammadpurwanto/e-commerce-golang/config" // Path ke config
	"github.com/Muhammadpurwanto/e-commerce-golang/router" // Path ke router
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Memuat file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 1. Inisialisasi Database
	db := config.InitDatabase()

	// 2. Inisialisasi Fiber App
	app := fiber.New()

	// 3. Setup Rute
	router.SetupRoutes(app, db) // Teruskan instance db ke router

	// 4. Menjalankan Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port %s\n", port)
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}