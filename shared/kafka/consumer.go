// pkg/kafka/consumer.go
package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageHandler is a function that processes a Kafka message
type MessageHandler func(ctx context.Context, message kafka.Message) error

// Consumer Kafka consumer wrapper
type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

// ConsumerConfig configuration for creating a consumer
type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config ConsumerConfig, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		MinBytes:       10e3,             // 10KB
		MaxBytes:       10e6,             // 10MB
		CommitInterval: time.Second,      // Her 1 saniyede commit
		StartOffset:    kafka.LastOffset, // En son mesajdan başla
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("Starting consumer for topic: %s, group: %s",
		c.reader.Config().Topic, c.reader.Config().GroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopped")
			return c.reader.Close()
		default:
			// Mesaj al
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			log.Printf("Received message: offset=%d key=%s", message.Offset, string(message.Key))

			// Mesajı işle
			if err := c.handler(ctx, message); err != nil {
				log.Printf("Error handling message: %v", err)
				continue // Hata olursa commit etme, tekrar işlensin
			}

			// Başarılı olursa commit et
			if err := c.reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
