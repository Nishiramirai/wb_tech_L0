package kafka

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

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
		Brokers:     cfg.Kafka.Brokers,
		Topic:       cfg.Kafka.Topic,
		GroupID:     "order-service",
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		service: svc,
		logger:  logger,
	}, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c.logger.Info("Kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Context cancelled, shutting down consumer...")
			return
		case sig := <-sigchan:
			c.logger.Info("Caught signal, shutting down...", zap.String("signal", sig.String()))
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
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

			var order model.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				c.logger.Error("Failed to unmarshal message", zap.Error(err), zap.ByteString("message_value", m.Value))
				continue
			}

			if err := c.service.SaveOrder(ctx, &order); err != nil {
				c.logger.Error("Failed to save order", zap.Error(err), zap.String("order_uid", order.OrderUID))
				continue
			}

		}
	}
}
