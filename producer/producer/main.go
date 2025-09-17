package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}
	// Конфигурация Writer с параметрами для надежной доставки
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_HOST")},
		Topic:    os.Getenv("TOPIC"),
		Balancer: &kafka.RoundRobin{},
	})
	defer writer.Close()

	// Чтение и отправка n JSON-файлов
	for i := 1; i <= 20; i++ {
		// Формирование пути к файлу (адаптируйте под вашу структуру файлов)
		filePath := filepath.Join("wildberies/L0/backend/Kafka/produce/files", fmt.Sprintf("file%d.json", i)) //поменять путь для докера (ъх по докеру)

		// Чтение файла
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Ошибка чтения файла %s: %v", filePath, err)
			continue
		}

		// Отправка сообщения в Kafka
		err = writer.WriteMessages(ctx, kafka.Message{
			Value: jsonData, // Отправка сырых JSON-данных
			// Key: []byte(fmt.Sprintf("key-%d", i)), // Опциональный ключ для гарантии порядка
			Time: time.Now(), // Метка времени
		})
		if err != nil {
			log.Printf("Ошибка отправки сообщения из файла %s: %v", filePath, err)
			// Здесь можно добавить логику повторной попытки или обработки ошибок
			continue
		}

		fmt.Printf("Сообщение из файла %s успешно отправлено\n", filePath)
	}

	fmt.Printf("Обработка 20 файлов завершена\n")
}
