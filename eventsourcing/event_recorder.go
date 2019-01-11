package eventsourcing

import (
	"github.com/coderbiq/dgo/model"
)

// EventRecorder record occured domain event
type EventRecorder struct {
	*model.EventRecorder
	sourced EventSourced
}

// EventRecorderFromSourced 创建一个支持 EventSourcing 的事件记录器
func EventRecorderFromSourced(sourced EventSourced, version uint) *EventRecorder {
	return &EventRecorder{
		EventRecorder: model.NewEventRecorder(version),
		sourced:       sourced,
	}
}

// RecordThan record an domain event
func (r *EventRecorder) RecordThan(event model.DomainEvent) {
	r.EventRecorder.RecordThan(event)
	if r.sourced != nil {
		r.sourced.Apply(event)
	}
}
