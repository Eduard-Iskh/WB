package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wildberies/L0/backend/internal/app"
	"wildberies/L0/backend/internal/domain"

	"github.com/go-chi/chi/v5"
)

type mockOrderService struct {
	createFn    func(ctx context.Context, order []byte) error
	getByIdFn   func(ctx context.Context, id string) (*domain.Order, error)
	warmCacheFn func(ctx context.Context, limit int) error
}

func (m *mockOrderService) Create(ctx context.Context, order []byte) error {
	if m.createFn != nil {
		return m.createFn(ctx, order)
	}
	return nil
}

func (m *mockOrderService) GetById(ctx context.Context, id string) (*domain.Order, error) {
	if m.getByIdFn != nil {
		return m.getByIdFn(ctx, id)
	}
	return nil, nil
}
func (m *mockOrderService) WarmCache(ctx context.Context, limit int) error {
	if m.warmCacheFn != nil {
		return m.warmCacheFn(ctx, limit)
	}
	return nil
}
func newApp(creatFn func(ctx context.Context, order []byte) error, getByIdFn func(ctx context.Context, id string) (*domain.Order, error), warmCacheFn func(ctx context.Context, limit int) error) *app.App {
	return &app.App{
		OrderService: &mockOrderService{
			createFn:    creatFn,
			getByIdFn:   getByIdFn,
			warmCacheFn: warmCacheFn,
		},
	}
}

func newOrder(id string) *domain.Order {
	return &domain.Order{
		OrderUID:          id,
		TrackNumber:       "TRACK-1",
		Entry:             "WBIL",
		Delivery:          domain.Delivery{Name: "Test"},
		Payment:           domain.Payment{Transaction: id},
		Items:             []domain.Item{{ChrtID: 1, TrackNumber: "TRACK-1", Price: 100, Rid: "rid-1", Name: "item", Size: "M", TotalPrice: 100, NmID: 1, Brand: "brand", Status: 1}},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer-1",
		DeliveryService:   "meest",
		Shardkey:          "1",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func TestGetOrderReturnServerFalls(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("service error")

	application := newApp(nil, func(ctx context.Context, id string) (*domain.Order, error) {
		return nil, wantErr
	}, nil)

	router := chi.NewRouter()
	router.Get("/order/{id}", GetOrder(application.OrderService))

	req := httptest.NewRequest(http.MethodGet, "/order/order-1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestGetORderReturnSuccess(t *testing.T) {
	t.Parallel()

	want := newOrder("order-1")

	application := newApp(nil,
		func(ctx context.Context, id string) (*domain.Order, error) {
			return want, nil
		},
		nil,
	)

	router := chi.NewRouter()
	router.Get("/order/{id}", GetOrder(application.OrderService))

	req := httptest.NewRequest(http.MethodGet, "/order/order-1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var response struct {
		Status string `json:"status"`
		Data   struct {
			Domain domain.Order `json:"order"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Fatalf("response status = %q, want %q", response.Status, "success")
	}

	if response.Data.Domain.OrderUID != want.OrderUID {
		t.Fatalf("response order uid = %q, want %q", response.Data.Domain.OrderUID, want.OrderUID)
	}
}
