package main

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "test-topic"
	broker := "kafka:9092"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic: topic,
	})

	defer writer.Close()

	for i := 0; i < 5; i++ {
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("key"),
				Value: []byte("Hello from Go producer!"),
			},
		)
		if err == nil {
			log.Println("Message sent successfully")
			break
		}
		log.Printf("Attempt %d failed: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
}