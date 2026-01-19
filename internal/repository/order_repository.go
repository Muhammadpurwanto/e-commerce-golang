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
	UpdateStatus(orderID uint, status string) (model.Order, error)
	FindOrderByIDForAdmin(orderID uint) (model.Order, error)
	CreateOrderFromCart(userID uint, cart model.Cart) (*model.Order, error)
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

func (r *orderRepository) UpdateStatus(orderID uint, status string) (model.Order, error) {
	var order model.Order
	err := r.db.Model(&order).Where("id = ?", orderID).Update("status", status).Error
	if err != nil {
		return model.Order{}, err
	}
	// Ambil data order terbaru setelah update
	r.db.Preload("User").Preload("OrderItems.Product").First(&order, orderID)
	return order, nil
}

// FindOrderByIDForAdmin tidak dibatasi oleh userID
func (r *orderRepository) FindOrderByIDForAdmin(orderID uint) (model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").Preload("OrderItems.Product").First(&order, orderID).Error
	return order, err
}

func (r *orderRepository) CreateOrderFromCart(userID uint, cart model.Cart) (*model.Order, error) {
	var createdOrder *model.Order

	// Mulai transaksi
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Hitung total amount dan siapkan order items
		var totalAmount float64
		var orderItems []model.OrderItem
		var productIDs []uint
		for _, cartItem := range cart.CartItems {
			productIDs = append(productIDs, cartItem.ProductID)
		}
		
		// Kunci baris produk yang relevan untuk mencegah race condition
		var products []model.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", productIDs).Find(&products).Error; err != nil {
			return err
		}

		productMap := make(map[uint]model.Product)
		for _, p := range products {
			productMap[p.ID] = p
		}

		for _, cartItem := range cart.CartItems {
			product, ok := productMap[cartItem.ProductID]
			if !ok {
				return errors.New("product not found")
			}
			
			// 2. Validasi stok
			if product.Stock < cartItem.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}

			itemTotal := product.Price * float64(cartItem.Quantity)
			totalAmount += itemTotal

			orderItems = append(orderItems, model.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     product.Price,
			})
		}
		
		// 3. Buat record Order utama
		order := &model.Order{
			UserID:      userID,
			TotalAmount: totalAmount,
			Status:      "pending",
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 4. Buat record OrderItem dan kurangi stok produk
		for i, item := range orderItems {
			// Update stok
			product := productMap[item.ProductID]
			newStock := product.Stock - item.Quantity
			if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", newStock).Error; err != nil {
				return err
			}
			// Tautkan OrderItem ke Order
			orderItems[i].OrderID = order.ID
		}
		
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// 5. Kosongkan keranjang: hapus semua CartItem dan Cart itu sendiri
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&cart).Error; err != nil {
			return err
		}

		createdOrder = order // Simpan order yang berhasil dibuat untuk dikembalikan
		return nil // Commit transaksi
	})

	if err != nil {
		return nil, err // Jika ada error, transaksi akan di-rollback secara otomatis
	}
    
    // Muat relasi untuk response yang lengkap
    r.db.Preload("User").Preload("OrderItems.Product").First(&createdOrder, createdOrder.ID)
	return createdOrder, nil
}