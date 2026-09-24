package order

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
	"wildberies/L0/backend/internal/domain"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
)

func cleanupTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := db.Exec(ctx, `
	TRUNCATE TABLE items, payments, deliveries, orders RESTART IDENTITY CASCADE
`)
	if err != nil {
		t.Fatalf("cleanup tables: %v", err)
	}

}

func openTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		t.Skip("TEST_DB_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		t.Fatalf("ping test database: %v", err)
	}
	t.Cleanup(pool.Close)

	return pool
}

func newTestOrder(id string, createdAt time.Time) *domain.Order {
	return &domain.Order{
		OrderUID:    id,
		TrackNumber: "TRAC" + id,
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: domain.Payment{
			Transaction:  id,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "TRACK-" + id,
				Price:       453,
				Rid:         "rid-" + id,
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:          "en",
		CustomerID:      "test",
		DeliveryService: "meest",
		Shardkey:        "9",
		SmID:            99,
		DateCreated:     createdAt,
		OofShard:        "1",
	}
}

func TestCreateAndGetByID(t *testing.T) {
	pool := openTestDB(t)
	cleanupTables(t, pool)

	ctx := context.Background()
	repo := NewOrderRepository(pool)

	createdAt := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	order := newTestOrder("order-1", createdAt)

	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	got, err := repo.GetById(ctx, "order-1")

	if err != nil {
		t.Fatalf("GetById() unexpected error: %v", err)
	}

	if diff := cmp.Diff(order, got); diff != "" {
		t.Fatalf("GetById() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetByIDReturnsNotFound(t *testing.T) {
	pool := openTestDB(t)
	cleanupTables(t, pool)

	ctx := context.Background()
	repo := NewOrderRepository(pool)

	got, err := repo.GetById(ctx, "order-1")

	if got != nil {
		t.Fatalf("GetByID() order = %v, want nil", got)
	}

	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("GetById() error = %v, want %v", err, domain.ErrOrderNotFound)
	}
}
