package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/ansonmathai/log-intel-system/internal/simulate"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "kafka:9092")
	topic := getEnv("KAFKA_TOPIC", "logs")
	mode := getEnv("PRODUCER_MODE", "stream") // "burst" or "stream"
	count := getIntEnv("LOG_COUNT", 10)
	interval := getIntEnv("LOG_INTERVAL_MS", 1000)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer writer.Close()

	log.Printf("🚀 Producer started in %s mode. Sending to topic: %s", mode, topic)

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	switch mode {
	case "burst":
		for i := 0; i < count; i++ {
			sendLog(writer)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	case "stream":
	loop:
		for {
			select {
			case <-stop:
				log.Println("🛑 Shutdown signal received. Exiting producer.")
				break loop
			default:
				sendLog(writer)
				time.Sleep(time.Duration(interval) * time.Millisecond)
			}
		}
	default:
		log.Fatalf("❌ Unknown PRODUCER_MODE: %s", mode)
	}
}

func sendLog(writer *kafka.Writer) {
	entry := simulate.GenerateLog()
	payload := simulate.FormatLog(entry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(entry.Category),
		Value: []byte(payload),
	})
	cancel()

	if err != nil {
		log.Printf("❌ Failed to send: %v", err)
	} else {
		log.Printf("✅ Sent: %s", payload)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}