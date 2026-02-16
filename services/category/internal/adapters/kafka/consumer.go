package kafka

import (
	"context"

	"product-app/shared/kafka"
)

type ConsumerAdapter struct {
	consumer *kafka.Consumer
}

func NewConsumerAdapter(brokers []string, topic, groupID string, handler kafka.MessageHandler) *ConsumerAdapter {
	return &ConsumerAdapter{
		consumer: kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}, handler),
	}
}

func (c *ConsumerAdapter) Start(ctx context.Context) error {
	return c.consumer.Start(ctx)
}

func (c *ConsumerAdapter) Close() error {
	return c.consumer.Close()
}
