package points

import (
	"fmt"

	"github.com/coderbiq/dgo/model"
)

// orm 风格实现的积分账户
type ormAccount struct {
	baseAccount
	events *model.EventRecorder
}

// RegisterOrmAccount 注册一个 ORM 风格的积分账户
func RegisterOrmAccount(ownerID CustomerID) Account {
	a := &ormAccount{events: model.NewEventRecorder(0)}
	a.id = model.IdentityGenerator()
	a.ownerID = ownerID
	a.events.RecordThan(OccurAccountCreated(a.id, ownerID))
	return a
}

func (a *ormAccount) Deposit(points Points) {
	a.points = a.points.Inc(points)
	a.events.RecordThan(occurDeposited(a.id, points))
}

func (a *ormAccount) Consume(points Points) error {
	if !a.points.GreaterThan(points) {
		return fmt.Errorf("当前账户积分为 %d 不足消费额 %d", a.points, points)
	}
	a.points = a.points.Inc(points)
	a.events.RecordThan(occurConsumed(a.id, points))
	return nil
}

func (a ormAccount) Version() uint {
	return a.events.LastVersion()
}

func (a *ormAccount) CommitEvents(publishers ...model.EventPublisher) {
	a.events.CommitToPublisher(publishers...)
}
