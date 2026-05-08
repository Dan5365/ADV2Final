package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/final/notification-service/internal/consumer"
)

func main() {
	rabbitURL := getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbitmq dial: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("rabbitmq channel: %v", err)
	}
	defer ch.Close()

	c := consumer.New(ch)
	if err := c.Start(); err != nil {
		log.Fatalf("consumer: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
