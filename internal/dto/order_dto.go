package dto

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	OrderItems []OrderItemRequest `json:"order_items" validate:"required,min=1,dive"`
}