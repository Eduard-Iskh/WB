package services

import (
	"context"
	"log/slog"
	"time"
	"wildberies/L0/backend/cache"
	domain "wildberies/L0/backend/internal/entify"
	valid "wildberies/L0/backend/internal/services/validate"
)

type OrderService struct {
	orderRepo domain.OrderRepository
	logger    *slog.Logger
	cache     *cache.Cache
}

func NewOrderService(orderRepo domain.OrderRepository, logger *slog.Logger, cache *cache.Cache) domain.OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		logger:    logger,
		cache:     cache,
	}
}

func (r *OrderService) Create(ctx context.Context, order []byte) error {
	// Проверка валидности данных
	orderData, err := valid.ProcessValid(order)
	if err != nil {
		r.logger.Error("Data validation error ", slog.Any("error", err))
		return err
	}
	r.logger.Info("creating order:", "customer id", orderData.CustomerID)

	// Внесение новых данных в БД
	err = r.orderRepo.Create(ctx, orderData)
	if err != nil {
		r.logger.Error("Create new order error: ", slog.Any("error", err))
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
		r.logger.Info("Ошибка Get By ID:",
			slog.Any("error", err),
			slog.String("id", id))
		return nil, err
	}

	// Сохраняем в кэш для будущих запросов
	r.cache.Set(id, *order)

	elapsed := time.Since(start)
	r.logger.Info("Данные получены из БД и сохранены в кэш",
		slog.String("id", id),
		slog.Int("duration,", int(elapsed.Microseconds())))
	return order, nil
}
