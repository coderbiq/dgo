package points

import (
	"fmt"
	"time"

	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

type sourcedAccount struct {
	baseAccount
	events *eventsourcing.EventRecorder

	// 展示业务模型
	DepositedPoints Points    `json:"depositedPoints"`
	ConsumedPoints  Points    `json:"consumedPoints"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// RegisterSourcedAccount 注册一个 EventSourcing 非 CQRS 风格的积分账户
func RegisterSourcedAccount(ownerID model.StringID) Account {
	a := new(sourcedAccount)
	a.events = eventsourcing.EventRecorderFromSourced(a, 0)
	a.events.RecordThan(OccurAccountCreated(
		model.IDGenerator.LongID(),
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
		a.id = event.AggregateID().(model.LongID)
		a.ownerID = event.(AccountCreated).OwnerID()
		a.DepositedPoints = Points(0)
		a.ConsumedPoints = Points(0)
		a.CreatedAt = event.CreatedAt()
		a.UpdatedAt = event.CreatedAt()
		break
	case AccountDepositedEvent:
		points := event.(AccountDeposited).Points()
		a.points = a.points.Inc(points)
		a.DepositedPoints = a.DepositedPoints.Inc(points)
		a.UpdatedAt = event.CreatedAt()
		break
	case AccountConsumedEvent:
		points := event.(AccountConsumed).Points()
		a.points = a.points.Dec(points)
		a.ConsumedPoints = a.ConsumedPoints.Inc(points)
		a.UpdatedAt = event.CreatedAt()
		break
	}
}
