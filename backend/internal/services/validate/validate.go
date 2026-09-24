package valid

import (
	"encoding/json"
	"fmt"
	domain "wildberies/L0/backend/internal/domain"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ProcessValid(message []byte) (*domain.Order, error) {
	// Быстрая проверка валидности JSON
	if len(message) == 0 {
		return nil, fmt.Errorf("%w: empty order message", domain.ErrInvalidOrder)
	}

	// Парсинг JSON
	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		return nil, fmt.Errorf("%w: decode order json: %v ", domain.ErrInvalidOrder, err)
	}

	// Валидация структуры
	if err := validate.Struct(&order); err != nil {
		return nil, fmt.Errorf("%w: validate order: %v", domain.ErrInvalidOrder, err)
	}

	return &order, nil
}
