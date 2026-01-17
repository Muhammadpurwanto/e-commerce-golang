package handler

import (
	"strconv"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/dto"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo, productRepo: productRepo}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(float64)

	var orderItems []model.OrderItem
	var totalAmount float64

	for _, itemReq := range req.OrderItems {
		product, err := h.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found: ID " + strconv.Itoa(int(itemReq.ProductID))})
		}
		if product.Stock < itemReq.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient stock for product: " + product.Name})
		}
		
		itemTotal := product.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, model.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
		})
	}

	order := &model.Order{
		UserID:      uint(userID),
		TotalAmount: totalAmount,
		Status:      "pending",
	}

	createdOrder, err := h.orderRepo.Create(order, orderItems)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdOrder)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	orders, err := h.orderRepo.FindUserOrders(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	
	userID := c.Locals("user_id").(float64)
	
	order, err := h.orderRepo.FindOrderByID(uint(orderID), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}
	
	return c.JSON(order)
}