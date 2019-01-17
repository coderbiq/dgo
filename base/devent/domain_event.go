package devent

import (
	"context"
	"errors"
	"time"

	"github.com/coderbiq/dgo/base/vo"
)

// Event 定义领域事件外观
type Event interface {
	ID() vo.Identity
	Name() string
	AggregateID() vo.Identity
	Version() uint
	CreatedAt() time.Time
}

// Publisher 定义事件发布器
//
// 事件发布器可以将领域事件发布到消息中间件或事件存储，具体如何发布由发布器内部进行实现。
type Publisher interface {
	// Publish 发布领域事件
	//
	// 发布器内部应该自行维护失败状态。例如：提交消息中间件失败后可以将需要发布的事件暂存，待检查
	// 到消息中间件恢复后再继续发布。
	Publish(events ...Event)
}

// Producer 定义领域事件生产者
type Producer interface {
	CommitEvents(publishers ...Publisher)
}

// Consumer 定义领域事件消费者
type Consumer interface {
	Handle(Event)
}

// Router 根据事件名称返回订阅了该事件的所有消费者
type Router interface {
	Consumers(eventName string) ([]Consumer, bool)
}

// Bus 定义消息总线
type Bus interface {
	Publisher
	AddRouter(Router)
}

// ConsumerFunc 包装一个领域事件处理函数为领域事件消费者
func ConsumerFunc(handle func(Event)) Consumer {
	return &simpleConsumer{handle: handle}
}

type simpleConsumer struct {
	handle func(Event)
}

func (consume *simpleConsumer) Handle(event Event) {
	consume.handle(event)
}

// ValidEvent 验证一个领域事件是否完整
func ValidEvent(event Event) error {
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

type simpleBus struct {
	concurrent    uint
	eventChannel  chan Event
	handleChannel chan *simpleBusHandle
	routers       []Router
}

type simpleBusHandle struct {
	event    Event
	consumer Consumer
}

// SimpleBus 创建一个简单的消息总线
func SimpleBus(concurrent uint) Bus {
	return &simpleBus{
		concurrent:    concurrent,
		eventChannel:  make(chan Event, 1000),
		handleChannel: make(chan *simpleBusHandle),
		routers:       []Router{},
	}
}

func (bus *simpleBus) AddRouter(router Router) {
	bus.routers = append(bus.routers, router)
}

func (bus simpleBus) Publish(events ...Event) {
	for _, event := range events {
		bus.eventChannel <- event
	}
}

// TODO: context done 时如果事件流中有未处理的消息将丢失，应该把未处理消息持久化再退出
func (bus simpleBus) Run(ctx context.Context) {
	for i := uint(0); i < bus.concurrent; i++ {
		go func(ctx context.Context, c <-chan *simpleBusHandle) {
			for {
				select {
				case <-ctx.Done():
					return
				case handleInfo := <-c:
					handleInfo.consumer.Handle(handleInfo.event)
					break
				}
			}
		}(ctx, bus.handleChannel)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-bus.eventChannel:
			for _, router := range bus.routers {
				if consumers, has := router.Consumers(event.Name()); has {
					for _, consumer := range consumers {
						bus.handleChannel <- &simpleBusHandle{
							consumer: consumer,
							event:    event,
						}
					}
				}
			}
			break
		}
	}
}

type simpleRouter struct {
	consumers map[string][]Consumer
}

// SimpleRouter 创建一个简单的事件路由
func SimpleRouter(routes map[string][]Consumer) Router {
	router := &simpleRouter{consumers: routes}
	return router
}

func (router *simpleRouter) Consumers(eventName string) ([]Consumer, bool) {
	consumers, has := router.consumers[eventName]
	return consumers, has
}
