package points

import (
	"fmt"

	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

type sourcedCqrsAccount struct {
	baseAccount
	events *eventsourcing.EventRecorder
}

// RegisterSourcedCqrsAccount 注册一个 EventSourcing 风格的积分账户
func RegisterSourcedCqrsAccount(ownerID model.StringID) Account {
	a := new(sourcedCqrsAccount)
	a.events = eventsourcing.EventRecorderFromSourced(a, 0)
	a.events.RecordThan(OccurAccountCreated(
		model.IDGenerator.LongID(),
		ownerID))
	return a
}

func (a *sourcedCqrsAccount) Deposit(points Points) {
	a.events.RecordThan(occurDeposited(a.id, points))
}

func (a *sourcedCqrsAccount) Consume(points Points) error {
	if !a.points.GreaterThan(points) {
		return fmt.Errorf("当前账户积分为 %d 不足消费额 %d", a.points, points)
	}
	a.events.RecordThan(occurConsumed(a.id, points))
	return nil
}

func (a sourcedCqrsAccount) Version() uint {
	return a.events.LastVersion()
}

func (a *sourcedCqrsAccount) CommitEvents(publishers ...model.EventPublisher) {
	a.events.CommitToPublisher(publishers...)
}

func (a *sourcedCqrsAccount) Apply(event model.DomainEvent) {
	switch event.Name() {
	case AccountCreatedEvent:
		a.id = event.AggregateID().(model.LongID)
		a.ownerID = event.(AccountCreated).OwnerID()
		break
	case AccountDepositedEvent:
		a.points = a.points.Inc(event.(AccountDeposited).Points())
		break
	case AccountConsumedEvent:
		a.points = a.points.Dec(event.(AccountConsumed).Points())
		break
	}
}
