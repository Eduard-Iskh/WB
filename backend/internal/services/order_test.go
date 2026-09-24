package services

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"wildberies/L0/backend/internal/domain"
)

type mockOrderRepo struct {
	createFn     func(ctx context.Context, order *domain.Order) error
	getByIdFn    func(ctx context.Context, id string) (*domain.Order, error)
	getLatestFn  func(ctx context.Context, limit int) ([]*domain.Order, error)
	getByIdCalls int
	createCalls  int
}

func (m *mockOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	m.createCalls++
	if m.createFn != nil {
		return m.createFn(ctx, order)
	}
	return nil
}

func (m *mockOrderRepo) GetById(ctx context.Context, id string) (*domain.Order, error) {
	m.getByIdCalls++
	if m.getByIdFn != nil {
		return m.getByIdFn(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderRepo) GetLatest(ctx context.Context, limit int) ([]*domain.Order, error) {
	if m.getLatestFn != nil {
		return m.getLatestFn(ctx, limit)
	}
	return nil, nil
}

type mockOrderCache struct {
	data          map[string]domain.Order
	setChaceCalls int
}

func newMockOrderCache() *mockOrderCache {
	return &mockOrderCache{
		data: make(map[string]domain.Order),
	}
}
func (m *mockOrderCache) Set(id string, order domain.Order) {
	m.setChaceCalls++
	m.data[id] = order
}

func (m *mockOrderCache) Get(id string) (domain.Order, bool) {
	order, ok := m.data[id]
	return order, ok
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newTestOrderService(repo domain.OrderRepository, cache domain.OrderCache) *OrderService {
	return &OrderService{
		orderRepo: repo,
		cache:     cache,
		logger:    newTestLogger(),
	}

}

func testOrder(id string) *domain.Order {
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

func validOrderJSON() []byte {
	return []byte(`{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"name": "Test Testov",
			"phone": "+9720000000",
			"zip": "2639809",
			"city": "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region": "Kraiot",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "b563feb7b2b84b6test",
			"request_id": "",
			"currency": "USD",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "WBILMTESTTRACK",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			}
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`)
}

func TestOrderServiceGetByIdReturnOrderFromCache(t *testing.T) {
	t.Parallel()

	cache := newMockOrderCache()
	want := *testOrder("order-1")
	cache.Set(want.OrderUID, want)

	repo := &mockOrderRepo{
		getByIdFn: func(ctx context.Context, id string) (*domain.Order, error) {
			t.Fatal("repo не должен быть вызван при проверки cache")
			return nil, nil
		},
	}

	service := newTestOrderService(repo, cache)

	got, err := service.GetById(context.Background(), want.OrderUID)
	if err != nil {
		t.Fatalf("GetById() unexpected error: %v", err)
	}

	if got == nil {
		t.Fatal("GetById() returned nil order")
	}

	if got.OrderUID != want.OrderUID {
		t.Fatalf("GetById() order uid = %q, want %q", got.OrderUID, want.OrderUID)
	}

	if repo.getByIdCalls != 0 {
		t.Fatalf("repository GetById() calls = %d, want 0", repo.getByIdCalls)
	}
}

func TestOrderServiceGetByIdReturnOrderFromRep(t *testing.T) {
	t.Parallel()

	cache := newMockOrderCache()
	want := testOrder("order-2")

	repo := &mockOrderRepo{
		getByIdFn: func(ctx context.Context, id string) (*domain.Order, error) {
			return want, nil
		},
	}
	service := newTestOrderService(repo, cache)
	got, err := service.GetById(context.Background(), want.OrderUID)
	if err != nil {
		t.Fatalf("GetById() unexpected error: %v", err)
	}

	if got == nil {
		t.Fatalf("GetById() returned nil order")
	}

	if got.OrderUID != want.OrderUID {
		t.Fatalf("GetById() order uid = %q, want %q", got.OrderUID, want.OrderUID)
	}

	if repo.getByIdCalls != 1 {
		t.Fatalf("repository GetById() calls = %d, want 1", repo.getByIdCalls)
	}

	cached, ok := cache.Get(want.OrderUID)
	if !ok {
		t.Fatal("order was not saved to cache")
	}

	if cached.OrderUID != want.OrderUID {
		t.Fatalf("cached order uid = %q, want %q", cached.OrderUID, want.OrderUID)
	}
}

func TestOrderServiceGetByIdReturnErrFromRep(t *testing.T) {
	t.Parallel()

	cache := newMockOrderCache()
	wantErr := domain.ErrOrderNotFound

	repo := &mockOrderRepo{
		getByIdFn: func(ctx context.Context, id string) (*domain.Order, error) {
			return nil, wantErr
		},
	}

	service := newTestOrderService(repo, cache)

	got, err := service.GetById(context.Background(), "order-3")
	if err == nil {
		t.Fatal("GetById() error = nil, want error")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("GetById() error = %v, want %v", err, wantErr)
	}

	if got != nil {
		t.Fatalf("GetById() order = %v, want nil", got)
	}
}

func TestOrderServiceCreateReturnValidError(t *testing.T) {
	t.Parallel()

	repo := &mockOrderRepo{}
	cache := newMockOrderCache()

	service := newTestOrderService(repo, cache)

	invalidJSON := []byte(`{"order_uid": ""}`)
	err := service.Create(context.Background(), invalidJSON)
	if err == nil {
		t.Fatal("Create() error = nil, want validation error")
	}
	if repo.createCalls != 0 {
		t.Fatalf("repositore Create() calls = %d, want 0", repo.createCalls)
	}
	if cache.setChaceCalls != 0 {
		t.Fatalf("cache Set() calls = %d, want 0", cache.setChaceCalls)
	}

}

func TestOrderServiceCreateReturnRepError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("repository create error")

	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, order *domain.Order) error {
			return wantErr
		},
	}

	cache := newMockOrderCache()
	service := newTestOrderService(repo, cache)

	err := service.Create(context.Background(), validOrderJSON())
	if err == nil {
		t.Fatalf("Create() error = nil, want repository error")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("Create() error = %v, want %v", err, wantErr)
	}
}

func TestOrderServiceCreateStoresOrderSuccess(t *testing.T) {
	t.Parallel()

	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, order *domain.Order) error {
			return nil
		},
	}

	cache := newMockOrderCache()
	service := newTestOrderService(repo, cache)

	err := service.Create(context.Background(), validOrderJSON())

	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	order, ok := cache.Get("b563feb7b2b84b6test")

	if !ok {
		t.Fatalf("order was not stored in cache")
	}

	if order.OrderUID != "b563feb7b2b84b6test" {
		t.Fatalf("cached order uid = %q, want %q", order.OrderUID, "b563feb7b2b84b6test")
	}
	orderRep, err := service.GetById(context.Background(), "b563feb7b2b84b6test")
	if err != nil {
		t.Fatalf("GetById() unexpected error: %v", err)
	}
	if orderRep.OrderUID != "b563feb7b2b84b6test" {
		t.Fatalf("repository order uid = %q, want %q", orderRep.OrderUID, "b563feb7b2b84b6test")
	}
}

func TestOrderServiceWarmCache(t *testing.T) {
	t.Parallel()

	orders := []*domain.Order{
		testOrder("uid-1"),
		testOrder("uid-2"),
	}
	cache := newMockOrderCache()

	repo := &mockOrderRepo{
		getLatestFn: func(ctx context.Context, limit int) ([]*domain.Order, error) {
			if limit != 2 {
				t.Fatalf("GetLatest() limit = %d, want %d", limit, 2)
			}
			return orders, nil
		},
	}
	service := newTestOrderService(repo, cache)

	err := service.WarmCache(context.Background(), 2)

	if err != nil {
		t.Fatalf("WarmCache() unexpected error: %v", err)
	}

	for _, order := range orders {
		cacheOrder, ok := service.cache.Get(order.OrderUID)
		if !ok {
			t.Fatalf("order %q was not saved to cache", order.OrderUID)
		}

		if cacheOrder.OrderUID != order.OrderUID {
			t.Fatalf("cached order uid = %q, want %q", cacheOrder.OrderUID, order.OrderUID)
		}
	}
}
