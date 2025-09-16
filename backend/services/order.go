package services

import (
	"context"
	"encoding/json"

	"log/slog"
	"wildberies/L0/backend/cache"
	"wildberies/L0/backend/domain"
	valid "wildberies/L0/backend/validate"
)

type OrderService struct {
	orderRepo domain.OrderRepository
	logger    *slog.Logger
	cache     cache.Cache
}

func NewOrderService(orderRepo domain.OrderRepository, logger *slog.Logger) domain.OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (r *OrderService) Create(ctx context.Context, order *domain.Order) error {
	r.logger.Info("creating order:", "customer id", order.CustomerID) //

	data, _ := json.Marshal(order)
	// Проверка валидности данных
	orderData, err := valid.ProcessValid(data)
	if err != nil {
		r.logger.Error("Data validation error \n\n")
		return err
	}

	// Внесение новых данных в БД
	err = r.orderRepo.Create(ctx, orderData)
	if err != nil {
		r.logger.Error("Create new order error: \n\n")
		return err
	}

	// Внесение данных в cache
	r.cache.Set(orderData.OrderUID, *orderData)

	return nil
}

func (r *OrderService) GetById(ctx context.Context, id string) (*domain.Order, error) {
	order := &domain.Order{
		Delivery: domain.Delivery{},
		Payment:  domain.Payment{},
		Items:    []domain.Item{},
	}

	return order, nil
}
