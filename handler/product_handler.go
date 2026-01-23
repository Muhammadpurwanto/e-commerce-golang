package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/repository"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

// NewProductHandler membuat instance baru dari ProductHandler
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetAllProducts handler untuk mendapatkan semua produk
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	// Generate pagination dari query params
	pagination := utils.GeneratePaginationFromRequest(c)
	products, err := h.repo.FindAll(&pagination) // Kirim pagination ke repo

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// GetProductByID handler untuk mendapatkan produk berdasarkan ID
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

// CreateProduct handler untuk membuat produk baru
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product model.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	newProduct, err := h.repo.Create(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(newProduct)
}

// UpdateProduct handler untuk memperbarui produk
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Cek apakah produk ada
	existingProduct, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	var productData model.Product
	if err := c.BodyParser(&productData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	existingProduct.Name = productData.Name
	existingProduct.SKU = productData.SKU
	existingProduct.Price = productData.Price
	existingProduct.Stock = productData.Stock

	updatedProduct, err := h.repo.Update(existingProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedProduct)
}

// DeleteProduct handler untuk menghapus produk
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ProductHandler) UploadProductImage(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Cek apakah produk ada
	product, err := h.repo.FindByID(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image upload failed"})
	}

	// Buat nama file unik untuk menghindari konflik
	// Contoh: product-1-timestamp.jpg
	filename := fmt.Sprintf("product-%d-%d%s", product.ID, time.Now().Unix(), filepath.Ext(file.Filename))

	// Buat direktori 'uploads' jika belum ada
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}

	// Simpan file ke direktori ./uploads/
	err = c.SaveFile(file, "./uploads/"+filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Buat URL yang bisa diakses publik
	// Pastikan HOST dan PORT sesuai dengan environment Anda
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost"
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	imageURL := fmt.Sprintf("%s:%s/api/v1/uploads/%s", host, port, filename)

	// Update field ImageURL di database
	product.ImageURL = imageURL
	updatedProduct, err := h.repo.Update(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product with image URL"})
	}

	return c.JSON(updatedProduct)
}