package repository

import (
	"errors"

	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	Create(order *model.Order, items []model.OrderItem) (*model.Order, error)
	FindUserOrders(userID uint) ([]model.Order, error)
	FindOrderByID(orderID, userID uint) (model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order, items []model.OrderItem) (*model.Order, error) {
	// Memulai transaksi database
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat record Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Buat record OrderItem untuk setiap item
		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}

			// 3. Kurangi stok produk
			var product model.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, items[i].ProductID).Error; err != nil {
				return err
			}

			if product.Stock < items[i].Quantity {
				return errors.New("insufficient stock for product " + product.Name)
			}

			newStock := product.Stock - items[i].Quantity
			if err := tx.Model(&product).Update("stock", newStock).Error; err != nil {
				return err
			}
		}

		// Jika semua berhasil, commit transaksi
		return nil
	})

	if err != nil {
		return nil, err
	}
    
    // Load relasi User dan OrderItems untuk response
	r.db.Preload("User").Preload("OrderItems.Product").First(&order, order.ID)

	return order, nil
}

func (r *orderRepository) FindUserOrders(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindOrderByID(orderID, userID uint) (model.Order, error) {
	var order model.Order
	err := r.db.Preload("OrderItems.Product").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	return order, err
}