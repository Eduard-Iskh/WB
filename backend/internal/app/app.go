package app

import (
	"log/slog"
	"wildberies/L0/backend/domain"
	"wildberies/L0/backend/internal/config"
	order "wildberies/L0/backend/internal/storage"
	"wildberies/L0/backend/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config       *config.Config
	Logger       *slog.Logger
	OrderService domain.OrderRepository
}

func NewApp(db *pgxpool.Pool, logger *slog.Logger) *App {
	orderRepo := order.NewOrderRepository(db)

	return &App{
		Logger:       logger,
		OrderService: services.NewOrderService(orderRepo, logger),
	}
}
