package devent

import (
	"context"
)

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
func (bus *simpleBus) Run(ctx context.Context) {
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
