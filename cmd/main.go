package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muhammadpurwanto/e-commerce-golang/config"
	"github.com/Muhammadpurwanto/e-commerce-golang/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db := config.InitDatabase()
	app := fiber.New()
	router.SetupRoutes(app, db)

	// Channel untuk mendengarkan sinyal OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine untuk menjalankan server
	go func() {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Server is starting on port %s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Menunggu sinyal shutdown
	<-quit
	log.Println("Shutting down server...")

	// Memberi waktu 30 detik untuk menyelesaikan request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}