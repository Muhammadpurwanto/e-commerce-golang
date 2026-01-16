package handler

import "github.com/gofiber/fiber/v2"

// RootHandler menangani request ke root endpoint
func RootHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Selamat datang di E-commerce API v1!",
		"data":    nil,
	})
}