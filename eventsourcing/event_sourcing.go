package eventsourcing

import (
	"github.com/coderbiq/dgo/model"
)

// EventSourced 定义领域事件来源
// 实现了 EventSourcing 的聚合根，可以根据事件回放出聚合的状态。
type EventSourced interface {
	Apply(event model.DomainEvent)
}
