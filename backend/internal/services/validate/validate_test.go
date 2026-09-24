package valid

import (
	"errors"
	"testing"
	"time"
	domain "wildberies/L0/backend/internal/domain"

	"github.com/google/go-cmp/cmp"
)

func makeOrderJSON(itemsJSON string) []byte {
	return []byte(`{
		"order_uid": "order-test",
		"track_number": "TRACK-1",
		"entry": "WBIL",
		"delivery": {
			"name": "Test",
			"phone": "+7900",
			"zip": "01072003",
			"city": "Moscow",
			"address": "Ploshad Mira 15",
			"region": "Kraiot",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "order-test",
			"request_id": "",
			"currency": "USD",
			"provider": "wbpay",
			"amount": 1,
			"payment_dt": 123,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 333,
			"custom_fee": 0
		},
		"items": [` + itemsJSON +
		`],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "1",
		"sm_id": 1,
		"date_created": "2026-09-26T06:22:19Z",
		"oof_shard": "1"
	}`)
}

var emptyItemsJSON = makeOrderJSON("")
var emptyItemJSON = makeOrderJSON("{}")
var validAndEmptyItemJSON = makeOrderJSON(`
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
			}, {}`)
var validItemsJSON = makeOrderJSON(`
			{
				"chrt_id": 123,
				"track_number": "TRACK-1",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			},
			{
				"chrt_id": 234,
				"track_number": "TRACK-2",
				"price": 43,
				"rid": "ab4219099a764ae0btest",
				"name": "Mascaras",
				"sale": 40,
				"size": "1",
				"total_price": 357,
				"nm_id": 2389212,
				"brand": "Vivo",
				"status": 202
			}`)

func newExpectedOrder(items []domain.Item) *domain.Order {
	return &domain.Order{
		OrderUID:    "order-test",
		TrackNumber: "TRACK-1",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    "Test",
			Phone:   "+7900",
			Zip:     "01072003",
			City:    "Moscow",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: domain.Payment{
			Transaction:  "order-test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1,
			PaymentDT:    123,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   333,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "1",
		SmID:              1,
		DateCreated: time.Date(
			2026,
			time.September,
			26,
			6,
			22,
			19,
			0,
			time.UTC,
		),
		OofShard: "1",
	}
}
func TestProcessValid(t *testing.T) {
	expectedItems := []domain.Item{
		{
			ChrtID:      123,
			TrackNumber: "TRACK-1",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
		{
			ChrtID:      234,
			TrackNumber: "TRACK-2",
			Price:       43,
			Rid:         "ab4219099a764ae0btest",
			Name:        "Mascaras",
			Sale:        40,
			Size:        "1",
			TotalPrice:  357,
			NmID:        2389212,
			Brand:       "Vivo",
			Status:      202,
		},
	}
	testCases := []struct {
		name      string
		input     []byte
		wantOrder *domain.Order
		wantErr   error
	}{
		{
			name:    "empty items",
			input:   emptyItemsJSON,
			wantErr: domain.ErrInvalidOrder,
		},
		{
			name:    "item without required fields",
			input:   emptyItemJSON,
			wantErr: domain.ErrInvalidOrder,
		},
		{
			name:    "second item without required fields",
			input:   validAndEmptyItemJSON,
			wantErr: domain.ErrInvalidOrder,
		},
		{
			name:      "two valid items",
			input:     validItemsJSON,
			wantOrder: newExpectedOrder(expectedItems),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotOrder, err := ProcessValid(testCase.input)
			if testCase.wantErr != nil {
				if err == nil {
					t.Fatalf("ProcessValid() error = nil, want error matching %v", testCase.wantErr)
				}
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("ProcessValid() error = %v, want error matching %v", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ProcessValid() unexpected error: %v", err)
			}

			if diff := cmp.Diff(testCase.wantOrder, gotOrder); diff != "" {
				t.Fatalf("ProcessValid() result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
