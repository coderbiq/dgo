package model

// EventRecorder record occured domain event
type EventRecorder struct {
	version        uint
	recordedEvents []DomainEvent
}

// NewEventRecorder return new EventRecorder
func NewEventRecorder(version uint) *EventRecorder {
	return &EventRecorder{
		version:        version,
		recordedEvents: []DomainEvent{},
	}
}

// RecordThan record an domain event
func (r *EventRecorder) RecordThan(event DomainEvent) {
	r.version++
	r.recordedEvents = append(r.recordedEvents, event.WithVersin(r.version))
}

// RecordedEvents return domain events
func (r *EventRecorder) RecordedEvents() []DomainEvent {
	return r.recordedEvents
}

// LastVersion return last domain event version
func (r *EventRecorder) LastVersion() uint {
	return r.version
}
