package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig memuat konfigurasi dari file .env
func LoadConfig() (port string) {
	// Memuat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Mengambil nilai APP_PORT dari environment variable
	port = os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Port default jika tidak ada di .env
	}
	return port
}