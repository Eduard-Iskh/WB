package cache

import (
	"testing"
	"time"
	"wildberies/L0/backend/internal/domain"
)

func testDomainOrder(id string) domain.Order {
	return domain.Order{
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

func TestGet(t *testing.T) {
	test := []struct {
		name      string
		setup     func(c *Cache)
		id        string
		wantOrder string
		wantOk    bool
	}{
		{
			name: "return existing order",
			setup: func(c *Cache) {
				c.Set("order-1", testDomainOrder("order-1"))
			},
			id:        "order-1",
			wantOrder: "order-1",
			wantOk:    true,
		},
		{
			name:      "missing order",
			setup:     func(c *Cache) {},
			id:        "missing order",
			wantOrder: "",
			wantOk:    false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCache(2)
			tt.setup(c)

			got, ok := c.Get(tt.id)

			if ok != tt.wantOk {
				t.Fatalf("Get() ok = %v, want %v", ok, tt.wantOk)
			}

			if !tt.wantOk {
				return
			}

			if got.OrderUID != tt.wantOrder {
				t.Fatalf("Get() order uid = %q, want %q", got.OrderUID, tt.wantOrder)
			}
		})
	}
}

func TestNewCache(t *testing.T) {
	test := []struct {
		name         string
		maxItems     int
		wantMaxItems int
	}{
		{
			name:         "uses provided maxItems",
			maxItems:     5,
			wantMaxItems: 5,
		},

		{
			name:         "uses default when maxItems is zero",
			maxItems:     0,
			wantMaxItems: 100,
		},
		{
			name:         "uses default when maxItems is negative",
			maxItems:     -1,
			wantMaxItems: 100,
		},
	}

	for _, tt := range test {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCache(tt.maxItems)

			if got.maxItems != tt.wantMaxItems {
				t.Fatalf("maxItems = %d, want %d", got.maxItems, tt.wantMaxItems)
			}
		})

	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name          string
		maxItems      int
		initialOrders []domain.Order
		insertOrder   domain.Order
		wantPresent   []string
		wantMissing   []string
		wantTrackByID map[string]string
	}{
		{
			name:          "stores order in empty cache",
			maxItems:      2,
			initialOrders: nil,
			insertOrder:   testDomainOrder("order-1"),
			wantPresent:   []string{"order-1"},
			wantMissing:   nil,
			wantTrackByID: map[string]string{
				"order-1": "TRACK-1",
			},
		},
		{
			name:          "stores order without eviction when limit not reached",
			maxItems:      2,
			initialOrders: []domain.Order{testDomainOrder("order-1")},
			insertOrder:   testDomainOrder("order-2"),
			wantPresent:   []string{"order-1", "order-2"},
			wantMissing:   nil,
			wantTrackByID: map[string]string{
				"order-1": "TRACK-1",
				"order-2": "TRACK-1",
			},
		},
		{
			name:          "evicts oldest order when limit exceeded",
			maxItems:      2,
			initialOrders: []domain.Order{testDomainOrder("order-1"), testDomainOrder("order-2")},
			insertOrder:   testDomainOrder("order-3"),
			wantPresent:   []string{"order-2", "order-3"},
			wantMissing:   []string{"order-1"},
			wantTrackByID: map[string]string{
				"order-2": "TRACK-1",
				"order-3": "TRACK-1",
			},
		},
		{
			name:     "updates existing order without adding duplicate key",
			maxItems: 2,
			initialOrders: []domain.Order{
				testDomainOrder("order-1"),
				testDomainOrder("order-2"),
			},
			insertOrder: func() domain.Order {
				o := testDomainOrder("order-1")
				o.TrackNumber = "UPDATED-TRACK"
				return o
			}(),
			wantPresent: []string{"order-1", "order-2"},
			wantMissing: nil,
			wantTrackByID: map[string]string{
				"order-1": "UPDATED-TRACK",
				"order-2": "TRACK-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCache(tt.maxItems)

			for _, order := range tt.initialOrders {
				c.Set(order.OrderUID, order)
			}

			c.Set(tt.insertOrder.OrderUID, tt.insertOrder)

			for _, id := range tt.wantPresent {
				got, ok := c.Get(id)

				if !ok {
					t.Fatalf("Get(%q) ok = false, want true", id)
				}

				wantTrack := tt.wantTrackByID[id]
				if got.TrackNumber != wantTrack {
					t.Fatalf("Get(%q) track number = %q, want %q", id, got.TrackNumber, wantTrack)
				}

			}
			for _, id := range tt.wantMissing {
				if _, ok := c.Get(id); ok {
					t.Fatalf("Get(%q) ok = true, want false", id)
				}
			}

		})
	}
}
