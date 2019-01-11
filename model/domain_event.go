package model

import "time"

// Payload message payload
type Payload interface {
	JSON() ([]byte, error)
}

// DomainEvent domain event model
type DomainEvent interface {
	ID() Identity
	Name() string
	AggregateID() Identity
	Payload() Payload
	Version() uint
	CreatedAt() time.Time
	WithVersin(version uint) DomainEvent
}

// OccurDomainEvent create domain event
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
