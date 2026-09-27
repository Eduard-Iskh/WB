package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wildberies/L0/backend/internal/config"
	domain "wildberies/L0/backend/internal/domain"
	"wildberies/L0/backend/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type DLQMessage struct {
	Payload       []byte `json:"payload"`
	Error         string `json:"error"`
	OriginalTopic string `json:"original_topic"`
	Partition     int    `json:"partition"`
	Offset        int64  `json:"offset"`
}

func waitRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func ConsumerKafka(ctx context.Context, newOrder domain.OrderService, cfg config.Config, log *slog.Logger) error {

	// Создаем reader вместо DialLeader для непрерывного чтения
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaConfig.Brokers,
		Topic:   cfg.KafkaConfig.Topic,
		GroupID: cfg.KafkaConfig.GroupID,
	})
	dlqWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.KafkaConfig.Brokers,
		Topic:   cfg.KafkaConfig.DLQTopic,
	})
	defer dlqWriter.Close()
	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info("kafka consumer stopped", logger.Err(ctx.Err()))
				return nil
			}
			log.Error("read message error", logger.Err(err))
			return err
		}
		var processErr error

		for attempt := 1; attempt <= 3; attempt++ {

			processErr = newOrder.Create(ctx, msg.Value)
			if processErr == nil {
				break
			}

			if errors.Is(processErr, domain.ErrInvalidOrder) {
				log.Error(
					"process kafka order message failed",
					logger.Err(processErr),
					slog.Int("attempt", attempt))
				break
			}

			if errors.Is(processErr, domain.ErrOrderAlreadyExist) {
				log.Info(
					"process kafka order already exist",
					logger.Err(processErr),
					slog.Int("attempt", attempt))
				processErr = nil
				break
			}

			if attempt < 3 {
				if err = waitRetry(ctx, cfg.KafkaConfig.RetryDelay); err != nil {

					break
				}
			}
		}
		if ctx.Err() != nil {
			log.Info("kafka consumer stopped", logger.Err(ctx.Err()))
			return nil
		}
		if processErr != nil {
			dlqData := DLQMessage{
				Payload:       msg.Value,
				Error:         processErr.Error(),
				OriginalTopic: msg.Topic,
				Partition:     msg.Partition,
				Offset:        msg.Offset,
			}
			value, err := json.Marshal(dlqData)
			if err != nil {
				log.Error("process marshal DLQ messsage failed", logger.Err(err))
				return err
			}
			err = dlqWriter.WriteMessages(ctx, kafka.Message{
				Value: value,
				Key:   msg.Key,
			})
			if err != nil {
				log.Error("process DLQ writer failed", logger.Err(err))
				return err
			}
		}

		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Error("commit kafka message failed", logger.Err(err))
			return fmt.Errorf("commit message: %w", err)
		}
		// Подтверждаем обработку сообщения (опционально, зависит от настроек)
		log.Info("kafka message processed", slog.Int64("offset", msg.Offset))
	}
}
