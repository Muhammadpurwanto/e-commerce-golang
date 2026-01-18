package handler

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetProfile handler untuk mendapatkan profil user yang sedang login
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userID := c.Locals("user_id").(float64)

	user, err := h.userRepo.FindByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Jangan kirim password hash
	user.Password = ""

	return c.JSON(user)
}