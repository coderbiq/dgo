package model_test

import (
	"context"
	"fmt"

	"github.com/coderbiq/dgo/model"
)

type runner interface {
	Run(context.Context)
}

func ExampleSimpleEventBus() {

	eventBus := model.SimpleEventBus(5)
	go eventBus.(runner).Run(context.Background())

	eventBus.Listen(
		"accountCreated",
		model.EventConsumerFunc(func(event model.DomainEvent) {
			if e, ok := event.(*AccountCreated); ok {
				fmt.Printf(
					"account %s created, identity is %d\n",
					e.AccountName,
					e.AccountID.Int64(),
				)
				// Output:
				// account test created, identity is 1085182999932178432
				// account test2 created, identity is 1085182999932178433
			}
		}))

	eventBus.Publish(occurAccountCreate(model.IDGenerator.LongID(), "test"))
	eventBus.Publish(occurAccountCreate(model.IDGenerator.LongID(), "test2"))
}
