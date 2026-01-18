package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity uint `json:"quantity" validate:"required,gt=0"`
}