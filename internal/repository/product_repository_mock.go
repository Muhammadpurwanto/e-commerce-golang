package repository

import (
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/model"
	"github.com/Muhammadpurwanto/e-commerce-golang/internal/utils"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository adalah mock untuk ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindAll(pagination *utils.Pagination) ([]model.Product, error) {
	args := m.Called(pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(id uint) (model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == (model.Product{}) {
		return model.Product{}, args.Error(1)
	}
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductRepository) Create(product model.Product) (model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product model.Product) (model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}