package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"wildberies/L0/backend/cache"
	domain "wildberies/L0/backend/internal/entify"

	"github.com/segmentio/kafka-go"
)

func ConsumerKafka(ctx context.Context, newOrder domain.OrderService, cache *cache.Cache) {

	// Создаем reader вместо DialLeader для непрерывного чтения
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   os.Getenv("TOPIC"),
		GroupID: "order-service",
	})
	defer reader.Close()

	for n := 0; ; n++ {
		time.Sleep(2 * time.Second)

		select {
		case <-ctx.Done():
			log.Printf("Завершение работы Kafkla \n\n")
			// Срабатывает при отмене контекста
			// Корректное завершение работы
			return
		default:

			// Читаем сообщения в бесконечном цикле
			m, err := reader.ReadMessage(ctx)
			log.Printf("Получаем сообщение с offset №%d  из Kafka. \n", n)
			if err != nil {
				log.Printf("Read message error %v", err)
				continue
			}

			err = newOrder.Create(ctx, m.Value)
			if err != nil {
				log.Printf("Create new order error: %v", err)
				continue
			}

			// Подтверждаем обработку сообщения (опционально, зависит от настроек)
			fmt.Printf("Обработано сообщение c offset %d\n\n", m.Offset)
		}
	}
}
