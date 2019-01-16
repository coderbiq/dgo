package devent_test

import (
	"context"
	"fmt"

	"github.com/coderbiq/dgo/base/devent"
	"github.com/coderbiq/dgo/base/vo"
)

type runner interface {
	Run(context.Context)
}

func ExampleSimpleEventBus() {

	eventBus := devent.SimpleEventBus(5)
	go eventBus.(runner).Run(context.Background())

	eventBus.Listen(
		"accountCreated",
		devent.EventConsumerFunc(func(event devent.DomainEvent) {
			if e, ok := event.(*AccountCreated); ok {
				fmt.Printf(
					"account %s created, identity is %d\n",
					e.AccountName,
					e.AccountID.Int64(),
				)
			}
		}))

	eventBus.Publish(occurAccountCreate(vo.LongID(1), "test"))
	eventBus.Publish(occurAccountCreate(vo.LongID(2), "test2"))

	// Unordered output
	// account test created, identity is 1
	// account test2 created, identity is 2
}
