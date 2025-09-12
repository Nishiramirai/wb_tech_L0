package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"orders-service/internal/app/model"

	"github.com/segmentio/kafka-go"
)

func generateUID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b) + "test"
}

func main() {
	brokers := []string{"kafka:29092"}
	topic := "orders"

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	testOrder := model.Order{
		OrderUID:    generateUID(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{
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
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	orderJSON, err := json.Marshal(testOrder)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Value: orderJSON,
		},
	)
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	log.Printf("Successfully sent message with order UID: %s", testOrder.OrderUID)
}
