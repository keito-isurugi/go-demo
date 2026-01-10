package domain

import "time"

// DomainEvent ドメインイベントのインターフェース
// ドメインエキスパートが関心を持つ、ドメイン内で起きた出来事を表現する
type DomainEvent interface {
	// OccurredAt イベント発生日時を返す
	OccurredAt() time.Time
	// AggregateID 集約のIDを返す
	AggregateID() string
	// EventType イベントの種類を返す
	EventType() string
}
