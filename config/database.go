package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model" // Path ke model
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	// Ambil konfigurasi dari environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Membuat Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Membuka koneksi ke database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	// AutoMigrate untuk membuat tabel secara otomatis
	// Ini akan membuat tabel 'products' berdasarkan struct model.Product
	db.AutoMigrate(&model.Product{}, model.User{}, model.Order{}, model.OrderItem{}, model.Cart{}, model.CartItem{})
	log.Println("Database migrated.")

	return db
}