package dgo

// EventRecorder record occured domain event
type EventRecorder struct {
	version        uint
	recordedEvents []DomainEvent
	sourced        EventSourced
}

// NewEventRecorder return new EventRecorder
func NewEventRecorder(version uint) *EventRecorder {
	return &EventRecorder{
		version:        version,
		recordedEvents: []DomainEvent{},
	}
}

// EventRecorderFromSourced 创建一个支持 EventSourcing 的事件记录器
func EventRecorderFromSourced(sourced EventSourced, version uint) *EventRecorder {
	return &EventRecorder{
		version:        version,
		recordedEvents: []DomainEvent{},
		sourced:        sourced,
	}
}

// RecordThan record an domain event
func (r *EventRecorder) RecordThan(event DomainEvent) {
	r.version++
	r.recordedEvents = append(r.recordedEvents, event.WithVersin(r.version))
	if r.sourced != nil {
		r.sourced.Apply(event)
	}
}

// RecordedEvents return domain events
func (r *EventRecorder) RecordedEvents() []DomainEvent {
	return r.recordedEvents
}

// LastVersion return last domain event version
func (r *EventRecorder) LastVersion() uint {
	return r.version
}
