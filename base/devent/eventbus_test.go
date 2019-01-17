package devent_test

import (
	"context"
	"testing"
	"time"

	"github.com/coderbiq/dgo/base/devent"
	"github.com/coderbiq/dgo/base/vo"
)

func TestEventBus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	eventBus := devent.SimpleEventBus(5)
	go eventBus.(runner).Run(ctx)

	assert := false
	handleEvents := 0
	eventBus.AddRouter(devent.SimpleEventRouter(map[string][]devent.EventConsumer{
		"accountCreated": []devent.EventConsumer{
			devent.EventConsumerFunc(func(event devent.DomainEvent) {
				defer func() {
					if handleEvents == 2 {
						assert = true
						cancel()
					}
				}()
				if e, ok := event.(*AccountCreated); ok {
					if e.AccountName != "test" && e.AccountName != "test2" {
						cancel()
					}
					handleEvents++
				} else {
					cancel()
				}
			})},
	}))

	eventBus.Publish(occurAccountCreate(vo.IDGenerator.LongID(), "test"))
	eventBus.Publish(occurAccountCreate(vo.IDGenerator.LongID(), "test2"))

	<-ctx.Done()
	if !assert {
		t.Fail()
	}
}
