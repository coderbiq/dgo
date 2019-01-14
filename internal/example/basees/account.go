package basees

import (
	"time"

	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/internal/example/base"
	"github.com/coderbiq/dgo/internal/example/common"
	"github.com/coderbiq/dgo/model"
)

type account struct {
	common.BaseAccount
	base.AccountReadModel
	events *eventsourcing.EventRecorder
}

// RegisterAccount 为指定会员标识的会员注册一个新的积分账户
func RegisterAccount(ownerID model.StringID) base.Account {
	a := new(account)
	a.events = eventsourcing.EventRecorderFromSourced(a, 0)
	a.events.RecordThan(common.OccurAccountCreated(a.Identity, ownerID))
	return a
}

func (a account) Deposit(points common.Points) {
	a.events.RecordThan(common.OccurDeposited(a.Identity, points))
}

func (a account) Consume(points common.Points) error {
	if !a.CurPoints.GreaterThan(points) {
		return common.PointsInsufficientError{
			Current: a.CurPoints,
			Expect:  points,
		}
	}
	a.events.RecordThan(common.OccurConsumed(a.Identity, points))
	return nil
}

func (a *account) Apply(event model.DomainEvent) {
	switch event.Name() {
	case common.AccountCreatedEvent:
		a.Identity = event.AggregateID().(model.LongID)
		a.OwnerIdentity = event.(common.AccountCreated).OwnerID()
		a.DepPoints = common.Points(0)
		a.ConPoints = common.Points(0)
		a.Created = time.Now()
		a.Updated = time.Now()
		break
	case common.AccountDepositedEvent:
		points := event.(common.AccountDeposited).Points()
		a.CurPoints = a.CurPoints.Inc(points)
		a.DepPoints = a.DepPoints.Inc(points)
		a.Updated = event.CreatedAt()
		break
	case common.AccountConsumedEvent:
		points := event.(common.AccountConsumed).Points()
		a.CurPoints = a.CurPoints.Dec(points)
		a.ConPoints = a.ConPoints.Inc(points)
		a.Updated = time.Now()
		break
	}
}
