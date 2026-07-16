package responses

import (
	"backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	ShippingAddressID uuid.UUID `json:"shipping_address_id"`
	OrderNumber       string    `json:"order_number"`
	Status            string    `json:"status"`
	TotalAmount       float64   `json:"total_amount"`
	OrderedAt         time.Time `json:"ordered_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type OrderListResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	UserFullName      string    `json:"user_full_name"`
	Username          string    `json:"username"`
	UserEmail         string    `json:"user_email"`
	ShippingAddressID uuid.UUID `json:"shipping_address_id"`
	Phone             string    `json:"phone"`
	AddressLine       string    `json:"address_line"`
	City              string    `json:"city"`
	Province          string    `json:"province"`
	PostalCode        string    `json:"postal_code"`
	Country           string    `json:"country"`
	OrderNumber       string    `json:"order_number"`
	Status            string    `json:"status"`
	TotalAmount       float64   `json:"total_amount"`
	OrderedAt         time.Time `json:"ordered_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ToOrderResponse(order *models.Orders) *OrderResponse {
	return &OrderResponse{
		ID:                order.ID,
		UserID:            order.UserID,
		ShippingAddressID: order.ShippingAddressID,
		OrderNumber:       order.OrderNumber,
		Status:            order.Status,
		TotalAmount:       order.TotalAmount,
		OrderedAt:         order.OrderedAt,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}
}

func ToOrdersResponse(orders []models.Orders) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))

	for _, order := range orders {
		responses = append(responses, *ToOrderResponse(&order))
	}

	return responses
}
