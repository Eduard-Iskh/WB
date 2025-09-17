package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	// Создаем директорию для тестовых файлов
	os.MkdirAll("test_data", 0755)

	// Генерируем валидные файлы
	generateValidFiles()

	// Генерируем невалидные файлы
	generateInvalidFiles()

	fmt.Println("Сгенерировано 20 тестовых файлов в папке test_data")
}

func generateValidFiles() {
	// Базовый валидный заказ
	baseOrder := Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
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
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}

	// Генерируем 15 валидных вариантов
	for i := 1; i <= 15; i++ {
		order := baseOrder
		order.OrderUID = fmt.Sprintf("test_%d_%s", i, generateRandomString(10))

		// Вносим различные изменения для каждого заказа
		switch i {
		case 1:
			// Базовый валидный
		case 2:
			order.TrackNumber = fmt.Sprintf("TRACK_%d", i)
			order.Delivery.Name = fmt.Sprintf("Customer %d", i)
		case 3:
			// Добавляем второй товар
			order.Items = append(order.Items, Item{
				ChrtID:      9934931,
				TrackNumber: order.TrackNumber,
				Price:       200,
				Rid:         fmt.Sprintf("rid_%d", i),
				Name:        "Second Product",
				Sale:        10,
				Size:        "M",
				TotalPrice:  180,
				NmID:        2389213,
				Brand:       "Other Brand",
				Status:      200,
			})
		case 4:
			order.InternalSignature = "signature"
			order.Payment.RequestID = "req_123"
		case 5:
			order.Payment.Currency = "RUB"
		case 6:
			order.Payment.Provider = "yandex"
		case 7:
			order.Delivery.Region = "Moscow"
		case 8:
			order.DeliveryService = "dhl"
		case 9:
			order.Shardkey = "5"
			order.SmID = 50
		case 10:
			order.DateCreated = time.Now().Format(time.RFC3339)
		case 11:
			order.Items[0].Name = "Different Product"
			order.Items[0].Price = 300
		case 12:
			order.Delivery.Email = fmt.Sprintf("customer%d@example.com", i)
		case 13:
			order.Payment.Bank = "sberbank"
		case 14:
			order.Items[0].Status = 100
		case 15:
			order.Locale = "ru"
		}

		// Сохраняем в файл
		saveOrderToFile(order, fmt.Sprintf("test_data/valid_%d.json", i))
	}
}

func generateInvalidFiles() {
	// 1. Пустой файл
	os.WriteFile("test_data/invalid_empty.json", []byte(""), 0644)

	// 2. Отсутствует обязательное поле order_uid
	orderWithoutUID := Order{
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		// ... остальные поля, но без OrderUID
	}
	saveOrderToFile(orderWithoutUID, "test_data/invalid_missing_required.json")

	// 3. Неправильный тип данных
	invalidOrder := map[string]interface{}{
		"order_uid": "test_invalid_type",
		"payment": map[string]interface{}{
			"amount": "should_be_number_not_string", // Неправильный тип
		},
	}
	saveMapToFile(invalidOrder, "test_data/invalid_wrong_type.json")

	// 4. Дубликат (копия valid_1)
	data, _ := os.ReadFile("test_data/valid_1.json")
	var duplicate Order
	json.Unmarshal(data, &duplicate)
	duplicate.OrderUID = "b563feb7b2b84b6test" // Такой же order_uid как в базовом
	saveOrderToFile(duplicate, "test_data/invalid_duplicate_1.json")

	// 5. Невалидный email
	orderWithBadEmail := Order{
		OrderUID:    "test_bad_email",
		TrackNumber: "TESTTRACK",
		Entry:       "TEST",
		Delivery: Delivery{
			Name:    "Test User",
			Phone:   "+1234567890",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Test Address",
			Region:  "Moscow",
			Email:   "invalid-email", // Невалидный email
		},
		Payment: Payment{
			Transaction:  "test_bad_email",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "testpay",
			Amount:       1000,
			PaymentDT:    1637907727,
			Bank:         "test",
			DeliveryCost: 500,
			GoodsTotal:   500,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      111111,
				TrackNumber: "TESTTRACK",
				Price:       500,
				Rid:         "test_rid",
				Name:        "Test Product",
				Sale:        0,
				Size:        "M",
				TotalPrice:  500,
				NmID:        111111,
				Brand:       "Test Brand",
				Status:      200,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "test",
		Shardkey:          "1",
		SmID:              11,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}
	saveOrderToFile(orderWithBadEmail, "test_data/invalid_bad_email.json")
}

func saveOrderToFile(order Order, filename string) {
	jsonData, _ := json.MarshalIndent(order, "", "  ")
	os.WriteFile(filename, jsonData, 0644)
}

func saveMapToFile(data map[string]interface{}, filename string) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filename, jsonData, 0644)
}

func generateRandomString(length int) string {
	// Простая реализация для генерации случайной строки
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
