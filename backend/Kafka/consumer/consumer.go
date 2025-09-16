package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"wildberies/L0/backend/cache"
	"wildberies/L0/backend/domain"
	valid "wildberies/L0/backend/validate"

	"github.com/segmentio/kafka-go"
)

func ConsumerKafka(ctx context.Context, newOrder domain.OrderRepository, cache *cache.Cache) {

	// Создаем reader вместо DialLeader для непрерывного чтения
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   os.Getenv("TOPIC"),
		GroupID: "order-service",
	})
	defer reader.Close()

	for n := 0; ; n++ {
		time.Sleep(2 * time.Second)
		log.Printf("Получаем сообщение с offset № %d из Kafka. \n", n)
		select {
		case <-ctx.Done():
			// Срабатывает при отмене контекста
			// Корректное завершение работы
			return
		default:

			// Читаем сообщения в бесконечном цикле
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Read message error %v \n\n", err)
				continue
			}
			// service
			// Проверка валидности данных
			orderData, err := valid.ProcessValid(m.Value)
			if err != nil {
				log.Printf("Data validation error \n\n")
				continue
			}

			// Внесение данных в cache
			cache.Set(orderData.OrderUID, *orderData)

			// Внесение новых данных в БД
			err = newOrder.Create(ctx, orderData)
			if err != nil {
				log.Printf("Create new order error: %v \n\n", err)
				continue
			}

			// Подтверждаем обработку сообщения (опционально, зависит от настроек)
			fmt.Printf("Обработано сообщение c offset %d\n\n", m.Offset)
		}
	}
}
