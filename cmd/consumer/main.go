package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "logs"
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "log-consumer-group"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait: 5 * time.Second,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	log.Printf("📥 Consumer started. Listening to topic: %s", topic)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		msg, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			log.Printf("⚠️ Error reading message: %v", err)
			continue
		}

		fmt.Printf("📝 [%s] %s\n", msg.Time.Format(time.RFC3339), string(msg.Value))
	}
}