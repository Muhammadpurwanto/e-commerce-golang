package repository

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"gorm.io/gorm"
)

// ProductRepository mendefinisikan interface untuk operasi database produk
type ProductRepository interface {
	FindAll() ([]model.Product, error)
	FindByID(id uint) (model.Product, error)
	Create(product model.Product) (model.Product, error)
	Update(product model.Product) (model.Product, error)
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository membuat instance baru dari ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	return product, err
}

func (r *productRepository) Create(product model.Product) (model.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepository) Update(product model.Product) (model.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *productRepository) Delete(id uint) error {
	err := r.db.Delete(&model.Product{}, id).Error
	return err
}