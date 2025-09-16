package domain

import (
	"context"
)

type Order struct {
	OrderUID          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required,min=1"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	Shardkey          string   `json:"shardkey" validate:"required"`
	SmID              int      `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetById(ctx context.Context, id string) (*Order, error)
}

type OrderService interface {
	Create(ctx context.Context, order *Order) error
	GetById(ctx context.Context, id string) (*Order, error)
}

//type GetByIdInCache(ctx context.Context, id Order.OrderUID)

//getbyid(ctx, OrderUID)
