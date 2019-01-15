package model_test

import (
	"context"
	"testing"
	"time"

	"github.com/coderbiq/dgo/model"
)

func TestEventBus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	eventBus := model.SimpleEventBus(5)
	go eventBus.(runner).Run(ctx)

	assert := false
	handleEvents := 0
	eventBus.Listen(
		"accountCreated",
		model.EventConsumerFunc(func(event model.DomainEvent) {
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
		}))

	eventBus.Publish(occurAccountCreate(model.IDGenerator.LongID(), "test"))
	eventBus.Publish(occurAccountCreate(model.IDGenerator.LongID(), "test2"))

	<-ctx.Done()
	if !assert {
		t.Fail()
	}
}
