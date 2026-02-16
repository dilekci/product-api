package kafka

import (
	"context"

	"product-app/shared/kafka"
)

type ProducerAdapter struct {
	producer *kafka.Producer
}

func NewProducerAdapter(brokers []string, topic string) *ProducerAdapter {
	return &ProducerAdapter{
		producer: kafka.NewProducer(brokers, topic),
	}
}

func (p *ProducerAdapter) Publish(ctx context.Context, key string, value interface{}) error {
	return p.producer.PublishMessage(ctx, key, value)
}

func (p *ProducerAdapter) Close() error {
	return p.producer.Close()
}
