package points

import (
	"fmt"

	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

type sourcedAccount struct {
	baseAccount
	events *eventsourcing.EventRecorder
}

// RegisterSourcedAccount 注册一个 EventSourcing 风格的积分账户
func RegisterSourcedAccount(ownerID CustomerID) Account {
	a := new(sourcedAccount)
	a.events = eventsourcing.EventRecorderFromSourced(a, 0)
	a.events.RecordThan(OccurAccountCreated(
		model.IdentityGenerator(),
		ownerID))
	return a
}

func (a *sourcedAccount) Deposit(points Points) {
	a.events.RecordThan(occurDeposited(a.id, points))
}

func (a *sourcedAccount) Consume(points Points) error {
	if !a.points.GreaterThan(points) {
		return fmt.Errorf("当前账户积分为 %d 不足消费额 %d", a.points, points)
	}
	a.events.RecordThan(occurConsumed(a.id, points))
	return nil
}

func (a sourcedAccount) Version() uint {
	return a.events.LastVersion()
}

func (a *sourcedAccount) CommitEvents(publishers ...model.EventPublisher) {
	a.events.CommitToPublisher(publishers...)
}

func (a *sourcedAccount) Apply(event model.DomainEvent) {
	switch event.Name() {
	case AccountCreatedEvent:
		a.id = event.AggregateID()
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
