package consumer

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func New(ch *amqp.Channel) *Consumer {
	return &Consumer{channel: ch}
}

func (c *Consumer) Start() error {
	if err := c.channel.ExchangeDeclare("bookings", "topic", true, false, false, false, nil); err != nil {
		return err
	}

	q, err := c.channel.QueueDeclare("notifications", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for _, key := range []string{"booking.created", "booking.status_updated", "order.created", "order.status_updated"} {
		if err := c.channel.QueueBind(q.Name, key, "bookings", false, nil); err != nil {
			return err
		}
	}

	msgs, err := c.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("notification-service consuming messages...")
	for msg := range msgs {
		c.handle(msg)
	}
	return nil
}

func (c *Consumer) handle(msg amqp.Delivery) {
	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("unmarshal error: %v", err)
		msg.Nack(false, false)
		return
	}
	log.Printf("[notification] routing_key=%s payload=%v", msg.RoutingKey, payload)
	msg.Ack(false)
}
