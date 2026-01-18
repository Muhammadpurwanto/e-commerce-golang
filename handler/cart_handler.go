package handler

import (
	"strconv"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/dto"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartRepo repository.CartRepository
}

func NewCartHandler(cartRepo repository.CartRepository) *CartHandler {
	return &CartHandler{cartRepo: cartRepo}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	cart, err := h.cartRepo.GetCartByUserID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get cart"})
	}

	return c.JSON(cart)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	var req dto.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(float64)

	cartItem := model.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	cart, err := h.cartRepo.AddItemToCart(uint(userID), cartItem)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item to cart"})
	}

	return c.Status(fiber.StatusCreated).JSON(cart)
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	cartItemID, err := strconv.Atoi(c.Params("itemId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var req dto.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Otorisasi: Pastikan item ini milik user yang sedang login
	userID := c.Locals("user_id").(float64)
	item, err := h.cartRepo.GetCartItemByID(uint(cartItemID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
	}
	cart, _ := h.cartRepo.GetCartByUserID(uint(userID))
	if item.CartID != cart.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	updatedItem, err := h.cartRepo.UpdateCartItemQuantity(uint(cartItemID), req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cart item"})
	}

	return c.JSON(updatedItem)
}

func (h *CartHandler) RemoveCartItem(c *fiber.Ctx) error {
	cartItemID, err := strconv.Atoi(c.Params("itemId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	// Otorisasi
	userID := c.Locals("user_id").(float64)
	item, err := h.cartRepo.GetCartItemByID(uint(cartItemID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
	}
	cart, _ := h.cartRepo.GetCartByUserID(uint(userID))
	if item.CartID != cart.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	if err := h.cartRepo.RemoveCartItem(uint(cartItemID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove item from cart"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}