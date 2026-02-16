package ports

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, key string, value interface{}) error
	Close() error
}
