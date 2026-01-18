package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Pagination represents pagination info
type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

// GeneratePaginationFromRequest generates pagination info from request query
func GeneratePaginationFromRequest(c *fiber.Ctx) Pagination {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	sort := c.Query("sort", "created_at desc")

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

// Paginate returns a GORM scope for pagination
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(p.Limit).Order(p.Sort)
	}
}