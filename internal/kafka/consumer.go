package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"orders-service/internal/app/model"
	"orders-service/internal/app/service"
	"orders-service/internal/configs"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
	logger  *zap.Logger
}

func NewConsumer(cfg *configs.AppConfig, svc *service.OrderService, logger *zap.Logger) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		service: svc,
		logger:  logger,
	}, nil
}

func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka reader...")
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context) {
	c.logger.Info("Kafka consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Info("Context cancelled, shutting down consumer...")
				return
			}
			c.logger.Error("Failed to read message", zap.Error(err))
			continue
		}

		c.logger.Info("Received new message",
			zap.String("topic", m.Topic),
			zap.Int64("offset", m.Offset),
		)

		c.processMessage(ctx, m.Value)
	}
}

func (c *Consumer) processMessage(ctx context.Context, value []byte) {
	var order model.Order
	if err := json.Unmarshal(value, &order); err != nil {
		c.logger.Error("Failed to unmarshal message", zap.Error(err), zap.ByteString("message_value", value))
		return
	}

	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		err := c.service.SaveOrder(ctx, &order)
		if err == nil {
			return
		}

		c.logger.Warn("Failed to save order, retrying...",
			zap.Error(err),
			zap.String("order_uid", order.OrderUID),
			zap.Int("attempt", i+1),
		)

		select {
		case <-ctx.Done():
			c.logger.Warn("Context cancelled during retry loop")
			return
		case <-time.After(1 * time.Second):
		}
	}

	c.logger.Error("Failed to save order after multiple retries", zap.String("order_uid", order.OrderUID))
}
