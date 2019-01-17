package devent

import (
	"time"

	"github.com/coderbiq/dgo/base/vo"
)

// AggregateChanged 提供对 Event 基本外观的实现
type AggregateChanged struct {
	EventID          vo.LongID `json:"id"`
	EventName        string    `json:"name"`
	AggregateVersion uint      `json:"version"`
	ChangeTime       time.Time `json:"createdAt"`
}

// OccurAggregateChanged 组装一个聚合变更领域事件
func OccurAggregateChanged(name string, event Event) Event {
	if ac, ok := event.(aggregateChanged); ok {
		ac.withEventName(name)
	}
	return event
}

type aggregateChanged interface {
	init()
	withVersion(version uint)
	withEventName(name string)
}

// ID 返回聚合变更事件标识
func (ac AggregateChanged) ID() vo.Identity {
	return ac.EventID
}

// Name 返回领域事件名称
func (ac AggregateChanged) Name() string {
	return ac.EventName
}

// Version 返回聚合变更时的版本号
func (ac AggregateChanged) Version() uint {
	return ac.AggregateVersion
}

// CreatedAt 返回聚合变更的时间
func (ac AggregateChanged) CreatedAt() time.Time {
	return ac.ChangeTime
}

func (ac *AggregateChanged) init() {
	if ac.EventID.Empty() {
		ac.EventID = vo.IDGenerator.LongID()
	}
	if ac.ChangeTime.IsZero() {
		ac.ChangeTime = time.Now()
	}
}

func (ac *AggregateChanged) withVersion(version uint) {
	ac.AggregateVersion = version
}

func (ac *AggregateChanged) withEventName(name string) {
	ac.EventName = name
}
