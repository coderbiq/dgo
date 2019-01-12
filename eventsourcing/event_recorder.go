package eventsourcing

import (
	"github.com/coderbiq/dgo/model"
)

// EventRecorder 提供 EventSoucing 风格的事件记录器
//
// 要求使用本事件记录器的聚合需要实现 EventSourced 接口，与基础事件记录器的区别是在记录已发生
// 的领域事件后，会将事件应用到聚合通过事件生成聚合的最新状态。
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

// RecordThan 记录一个已发生的领域事件并将应该应用到发生事件的聚合
func (r *EventRecorder) RecordThan(event model.DomainEvent) {
	r.EventRecorder.RecordThan(event)
	if r.sourced != nil {
		r.sourced.Apply(event)
	}
}
