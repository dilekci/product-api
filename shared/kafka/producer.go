// pkg/kafka/producer.go
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer Kafka producer wrapper
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // Load balancing strategy

		// Performance tuning
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,

		// Compression
		Compression: kafka.Snappy,

		// Reliability
		RequiredAcks: kafka.RequireOne,
		Async:        false, // Synchronous mode for reliability

		// Retry configuration
		MaxAttempts:     3,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
	}

	return &Producer{writer: writer}
}

// PublishMessage sends a single message to Kafka
func (p *Producer) PublishMessage(ctx context.Context, key string, value interface{}) error {
	// Marshal value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
		Time:  time.Now(),
	}

	// Write message
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", p.writer.Topic, err)
	}

	return nil
}

// PublishMessageWithHeaders sends a message with custom headers
func (p *Producer) PublishMessageWithHeaders(ctx context.Context, key string, value interface{}, headers map[string]string) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Convert headers
	kafkaHeaders := make([]kafka.Header, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	message := kafka.Message{
		Key:     []byte(key),
		Value:   valueBytes,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// PublishBatch sends multiple messages in a single batch
func (p *Producer) PublishBatch(ctx context.Context, messages []kafka.Message) error {
	if len(messages) == 0 {
		return nil
	}

	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		return fmt.Errorf("failed to write batch messages: %w", err)
	}

	return nil
}

// PublishRawMessage sends a pre-constructed Kafka message
func (p *Producer) PublishRawMessage(ctx context.Context, message kafka.Message) error {
	err := p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write raw message: %w", err)
	}
	return nil
}

// Close closes the producer and releases resources
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// Stats returns producer statistics
func (p *Producer) Stats() kafka.WriterStats {
	return p.writer.Stats()
}
