package repository

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	GetCartByUserID(userID uint) (model.Cart, error)
	AddItemToCart(userID uint, item model.CartItem) (model.Cart, error)
	UpdateCartItemQuantity(cartItemID uint, quantity uint) (model.CartItem, error)
	RemoveCartItem(cartItemID uint) error
	GetCartItemByID(cartItemID uint) (model.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// GetCartByUserID mengambil atau membuat keranjang untuk user
func (r *cartRepository) GetCartByUserID(userID uint) (model.Cart, error) {
	var cart model.Cart
	// Preload untuk memuat relasi CartItems beserta Product di dalamnya
	err := r.db.Preload("CartItems.Product").FirstOrCreate(&cart, model.Cart{UserID: userID}).Error
	return cart, err
}

// AddItemToCart menambahkan item ke keranjang atau memperbarui quantity jika sudah ada
func (r *cartRepository) AddItemToCart(userID uint, item model.CartItem) (model.Cart, error) {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return model.Cart{}, err
	}

	item.CartID = cart.ID

	// Logika: Jika item produk sudah ada, update quantity. Jika tidak, buat baru.
	// OnConflict(...) DoUpdates(...) adalah fitur upsert dari GORM
	err = r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "cart_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
	}).Create(&item).Error

	if err != nil {
		return model.Cart{}, err
	}

	// Muat ulang data keranjang untuk mendapatkan state terbaru
	return r.GetCartByUserID(userID)
}

func (r *cartRepository) UpdateCartItemQuantity(cartItemID uint, quantity uint) (model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Model(&cartItem).Where("id = ?", cartItemID).Update("quantity", quantity).Error
	if err != nil {
		return model.CartItem{}, err
	}
	// Ambil data terbaru
	r.db.Preload("Product").First(&cartItem, cartItemID)
	return cartItem, nil
}

func (r *cartRepository) RemoveCartItem(cartItemID uint) error {
	return r.db.Delete(&model.CartItem{}, cartItemID).Error
}

func (r *cartRepository) GetCartItemByID(cartItemID uint) (model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.First(&cartItem, cartItemID).Error
	return cartItem, err
}