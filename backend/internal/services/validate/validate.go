package valid

import (
	"encoding/json"
	"errors"
	"log"
	domain "wildberies/L0/backend/internal/entify"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ProcessValid(message []byte) (*domain.Order, error) {
	// Быстрая проверка валидности JSON
	if len(message) == 0 {
		return nil, errors.New("ОШИБКА: пустой файл")
	}

	// Парсинг JSON
	var order domain.Order
	if err := json.Unmarshal(message, &order); err != nil {
		log.Printf("ОШИБКА: неправильная модель данных JSON ")
		return nil, err
	}

	// Валидация структуры
	if err := validate.Struct(&order); err != nil {
		log.Printf("ОШИБКА: есть пустые или заполненные неверно поля данных")
		return nil, err
	}

	return &order, nil
}
