package model

import (
	"errors"
	"time"
)

// DomainEvent 定义领域事件外观
type DomainEvent interface {
	ID() Identity
	Name() string
	AggregateID() Identity
	Version() uint
	CreatedAt() time.Time
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

// EventConsumer 定义领域事件消费者
type EventConsumer interface {
	Handle(DomainEvent)
}

// EventBus 定义消息总线
type EventBus interface {
	EventPublisher
	Listen(eventName string, consumer EventConsumer)
}

// EventConsumerFunc 包装一个领域事件处理函数为领域事件消费者
func EventConsumerFunc(handle func(DomainEvent)) EventConsumer {
	return &simpleEventConsumer{handle: handle}
}

type simpleEventConsumer struct {
	handle func(DomainEvent)
}

func (consume *simpleEventConsumer) Handle(event DomainEvent) {
	consume.handle(event)
}

// ValidDomainEvent 验证一个领域事件是否完整
func ValidDomainEvent(event DomainEvent) error {
	if event.ID() == nil || event.ID().Empty() {
		return errors.New("事件标识不能为空")
	}
	if event.Name() == "" {
		return errors.New("事件名称不能为空")
	}
	if event.AggregateID() == nil || event.AggregateID().Empty() {
		return errors.New("发生领域事件的聚合标识不能为空")
	}
	if event.Version() < 1 {
		return errors.New("无效的领域事件版本")
	}
	if event.CreatedAt().IsZero() {
		return errors.New("领域事件发生时间不能为空")
	}
	return nil
}

// AggregateChanged 提供对 DomainEvent 基本外观的实现
type AggregateChanged struct {
	EventID          LongID    `json:"id"`
	EventName        string    `json:"name"`
	AggregateVersion uint      `json:"version"`
	ChangeTime       time.Time `json:"createdAt"`
}

// OccurAggregateChanged 组装一个聚合变更领域事件
func OccurAggregateChanged(name string, event DomainEvent) DomainEvent {
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
func (ac AggregateChanged) ID() Identity {
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
		ac.EventID = IDGenerator.LongID()
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

type simpleEventBus struct {
	listeners map[string][]chan<- DomainEvent
}

// SimpleEventBus 创建一个简单的消息总线
func SimpleEventBus() EventBus {
	return &simpleEventBus{listeners: map[string][]chan<- DomainEvent{}}
}

func (bus simpleEventBus) Publish(events ...DomainEvent) {
	for _, event := range events {
		listeners, has := bus.listeners[event.Name()]
		if !has {
			continue
		}
		for _, listener := range listeners {
			go func(listener chan<- DomainEvent) {
				listener <- event
			}(listener)
		}
	}
}

func (bus simpleEventBus) Listen(eventName string, consumer EventConsumer) {
	c := make(chan DomainEvent, 10)
	go func(c <-chan DomainEvent) {
		for {
			select {
			case e := <-c:
				consumer.Handle(e)
			}
		}
	}(c)
	listeners, has := bus.listeners[eventName]
	if !has {
		listeners = []chan<- DomainEvent{}
	}
	listeners = append(listeners, c)
	bus.listeners[eventName] = listeners
}
