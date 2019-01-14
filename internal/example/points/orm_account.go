package points

import (
	"fmt"
	"time"

	"github.com/coderbiq/dgo/model"
)

// ORM 非 CQRS 风格实现的积分账户聚合
//
// 非 CQRS 风格的聚合因为没有读模型分担查询业务的需求，所以需要将展示业务和修改业务的能力
// 一起实现在聚合中。
type ormAccount struct {
	baseAccount
	events *model.EventRecorder

	// 展示业务模型
	DepositedPoints Points    `json:"depositedPoints"`
	ConsumedPoints  Points    `json:"consumedPoints"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// RegisterOrmAccount 注册一个 ORM 非 CQRS 风格的积分账户
func RegisterOrmAccount(ownerID model.StringID) Account {
	a := &ormAccount{events: model.NewEventRecorder(0)}
	a.id = model.IDGenerator.LongID()
	a.ownerID = ownerID
	a.DepositedPoints = Points(0)
	a.ConsumedPoints = Points(0)
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.events.RecordThan(OccurAccountCreated(a.id, ownerID))
	return a
}

func (a *ormAccount) Deposit(points Points) {
	a.points = a.points.Inc(points)
	a.DepositedPoints = a.DepositedPoints.Inc(points)
	a.UpdatedAt = time.Now()
	a.events.RecordThan(occurDeposited(a.id, points))
}

func (a *ormAccount) Consume(points Points) error {
	if !a.points.GreaterThan(points) {
		return fmt.Errorf("当前账户积分为 %d 不足消费额 %d", a.points, points)
	}
	a.points = a.points.Dec(points)
	a.ConsumedPoints = a.ConsumedPoints.Inc(points)
	a.UpdatedAt = time.Now()
	a.events.RecordThan(occurConsumed(a.id, points))
	return nil
}

func (a ormAccount) Version() uint {
	return a.events.LastVersion()
}

func (a *ormAccount) CommitEvents(publishers ...model.EventPublisher) {
	a.events.CommitToPublisher(publishers...)
}
