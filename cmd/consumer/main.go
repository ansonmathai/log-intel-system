package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := "kafka:9092"
	topic := "test-topic"
	groupID := "log-consumer-group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic: topic,
		GroupID: groupID,
		StartOffset: kafka.FirstOffset,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing reader: %v", err)
		}
	}()

	log.Println("Consumer started. Listening for messages...")
	
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		msg, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		fmt.Printf("Received: key=%s value=%s\n", string(msg.Key), string(msg.Value))
	}
}