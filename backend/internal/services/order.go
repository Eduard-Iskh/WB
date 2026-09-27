package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	domain "wildberies/L0/backend/internal/domain"
	valid "wildberies/L0/backend/internal/services/validate"
)

type OrderService struct {
	orderRepo domain.OrderRepository
	logger    *slog.Logger
	cache     domain.OrderCache
}

func NewOrderService(orderRepo domain.OrderRepository, cache domain.OrderCache, log *slog.Logger) domain.OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cache:     cache,
		logger:    log,
	}
}

func (r *OrderService) Create(ctx context.Context, order []byte) error {

	// Проверка валидности данных
	orderData, err := valid.ProcessValid(order)
	if err != nil {
		r.logger.Error("validate order message failed", slog.Any("error", err))
		return err
	}
	r.logger.Info("creating order", "customer id", orderData.CustomerID)

	// Внесение новых данных в БД
	err = r.orderRepo.Create(ctx, orderData)
	if err != nil {
		r.logger.Error("create order failed", slog.Any("error", err))
		return err
	}

	// Внесение данных в cache
	r.cache.Set(orderData.OrderUID, *orderData)

	return nil
}

func (r *OrderService) GetById(ctx context.Context, id string) (*domain.Order, error) {
	start := time.Now() // Начинаем отсчет времени

	// Сначала проверяем кэш
	if cachedOrder, exists := r.cache.Get(id); exists {
		elapsed := time.Since(start)
		r.logger.Info("Данные получены из кэша",
			slog.String("id", id),
			slog.Int("duration", int(elapsed.Nanoseconds())))
		return &cachedOrder, nil
	}

	// Если в кэше нет, получаем из репозитория
	order, err := r.orderRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			r.logger.Info("order not found",
				slog.Any("error", err),
				slog.String("id", id))
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	// Сохраняем в кэш для будущих запросов
	r.cache.Set(id, *order)

	elapsed := time.Since(start)
	r.logger.Info("Данные получены из БД и сохранены в кэш",
		slog.String("id", id),
		slog.Int("duration", int(elapsed.Microseconds())))
	return order, nil
}

func (r *OrderService) WarmCache(ctx context.Context, limit int) error {
	orders, err := r.orderRepo.GetLatest(ctx, limit)
	if err != nil {
		r.logger.Error("failed to warm cache", slog.Any("error", err))
		return fmt.Errorf("warm cache: %w", err)
	}
	for _, order := range orders {
		r.cache.Set(order.OrderUID, *order)
	}

	r.logger.Info("cache warmed", slog.Int("orders_count", len(orders)))
	return nil
}
