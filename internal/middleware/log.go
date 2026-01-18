package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger adalah middleware untuk logging request yang lebih detail
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Memberi ID unik untuk setiap request
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		start := time.Now()

		// Lanjutkan ke handler berikutnya
		err := c.Next()

		stop := time.Now()
		latency := stop.Sub(start)

		// Format log
		log.Printf(
			"[%s] %s | %s | %d | %13v | %s | %s",
			"INFO",
			requestID,
			c.Method(),
			c.Response().StatusCode(),
			latency,
			c.IP(),
			c.Path(),
		)

		return err
	}
}