package application

import (
	"context"
	"ddd/domain"
)

// EventPublisher ドメインイベントを発行するインターフェース
type EventPublisher interface {
	Publish(ctx context.Context, event domain.DomainEvent) error
}
