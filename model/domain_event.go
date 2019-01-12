package model

import (
	"time"
)

// DomainEvent 定义领域事件外观
type DomainEvent interface {
	ID() Identity
	Name() string
	AggregateID() Identity
	Payload() interface{}
	Version() uint
	CreatedAt() time.Time
	// WithVersion 使用提供的版本号产生一个新事件
	WithVersion(version uint) DomainEvent
}

// EventPublisher 定义事件发布器
//
// 事件发布器可以将领域事件发布到消息中间件或事件存储，具体如何发布由发布器内部进行实现。
type EventPublisher interface {
	// Publish 发布领域事件
	//
	// 发布器内部应该自行维护失败状态。例如：提交消息中间件失败后可以将需要发布的事件暂存，待检查
	// 到消息中间件恢复后再继续发布。
	Publish(events ...DomainEvent)
}

// EventProducer 定义领域事件生产者
type EventProducer interface {
	CommitEvents(publishers ...EventPublisher)
}

// OccurEvent 创建一个新发生的事件
func OccurEvent(aid Identity, name string, p interface{}) DomainEvent {
	return domainEvent{
		Message:     NewMessage(name, p),
		aggregateID: aid,
	}
}

type domainEvent struct {
	Message

	aggregateID Identity
	version     uint
}

func (e domainEvent) AggregateID() Identity {
	return e.aggregateID
}

func (e domainEvent) Version() uint {
	return e.version
}

func (e domainEvent) WithVersion(version uint) DomainEvent {
	return domainEvent{
		Message:     e.Message,
		aggregateID: e.aggregateID,
		version:     version,
	}
}
