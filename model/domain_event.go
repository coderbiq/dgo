package model

import "time"

// Payload 定义消息有效载荷外观
//
// 消息有效载荷可以应用于领域事件、CQRS指令中用于描述数据信息，为了将这些数据可以便于存储和在消息
// 中间件中传输需要为数据提交编码行为，这里统一定义为编码成 JSON 字符串。
type Payload interface {
	JSON() ([]byte, error)
}

// DomainEvent 定义领域事件外观
type DomainEvent interface {
	ID() Identity
	Name() string
	AggregateID() Identity
	Payload() Payload
	Version() uint
	CreatedAt() time.Time
	WithVersin(version uint) DomainEvent
}

// OccurDomainEvent 创建一个已经发生的领域事件
func OccurDomainEvent(aggregateID Identity,
	name string, payload Payload) DomainEvent {
	return baseDomainEvent{
		baseMessage: newBaseMessage(name, payload),
		aggregateID: aggregateID,
	}
}

type baseDomainEvent struct {
	baseMessage

	aggregateID Identity
	version     uint
}

func (e baseDomainEvent) ID() Identity {
	return e.id
}

func (e baseDomainEvent) Name() string {
	return e.name
}

func (e baseDomainEvent) AggregateID() Identity {
	return e.aggregateID
}

func (e baseDomainEvent) Payload() Payload {
	return e.payload
}

func (e baseDomainEvent) Version() uint {
	return e.version
}

func (e baseDomainEvent) CreatedAt() time.Time {
	return e.created
}

func (e baseDomainEvent) WithVersin(version uint) DomainEvent {
	return baseDomainEvent{
		baseMessage: e.baseMessage,
		aggregateID: e.aggregateID,
		version:     version,
	}
}
