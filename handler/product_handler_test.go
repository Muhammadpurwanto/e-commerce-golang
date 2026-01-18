package handler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestProductHandler_GetAllProducts(t *testing.T) {
	// Setup
	app := fiber.New()
	mockRepo := new(repository.MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	// Skenario 1: Sukses mendapatkan produk
	t.Run("Success", func(t *testing.T) {
		// Harapan: method FindAll di mockRepo akan dipanggil dan mengembalikan data
		expectedProducts := []model.Product{{ID: 1, Name: "Laptop"}}
		mockRepo.On("FindAll", &utils.Pagination{Limit: 10, Page: 1, Sort: "created_at desc"}).Return(expectedProducts, nil).Once()

		req := httptest.NewRequest("GET", "/api/v1/products", nil)
		app.Get("/api/v1/products", productHandler.GetAllProducts)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		// Anda bisa menambahkan asserstion untuk mengecek body response
		mockRepo.AssertExpectations(t) // Pastikan mock dipanggil
	})

	// Skenario 2: Gagal karena error dari repository
	t.Run("Repository Error", func(t *testing.T) {
		// Harapan: method FindAll akan dipanggil dan mengembalikan error
		mockRepo.On("FindAll", &utils.Pagination{Limit: 10, Page: 1, Sort: "created_at desc"}).Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest("GET", "/api/v1/products", nil)
		app.Get("/api/v1/products", productHandler.GetAllProducts)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})
}