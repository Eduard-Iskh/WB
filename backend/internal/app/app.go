package app

import (
	"log/slog"
	domain "wildberies/L0/backend/internal/domain"
	order "wildberies/L0/backend/internal/repository/postgres"
	"wildberies/L0/backend/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	OrderService domain.OrderService
}

func NewApp(db *pgxpool.Pool, cache domain.OrderCache, log *slog.Logger) *App {
	orderRepo := order.NewOrderRepository(db)

	return &App{
		OrderService: services.NewOrderService(orderRepo, cache, log),
	}
}
