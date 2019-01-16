package eventsourcing

import "github.com/coderbiq/dgo/base/devent"

// EventSourced 定义领域事件来源
//
// 领域事件来源是一个实现了 EventSourcing 的聚合，可以应用一组已经在聚合中发生的领域事件来构建
// 出聚合的状态。
type EventSourced interface {
	Apply(event devent.DomainEvent)
}
