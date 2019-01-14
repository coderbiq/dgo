package points

import (
	"fmt"

	"github.com/coderbiq/dgo/model"
)

// ORM & CQRS 风格实现的积分账户聚合
//
// CQRS 风格的聚合只需要关心写操作相关的模型特征，不用关心应用复杂展示场景下的展示数据需求。所
// 以不需要包含例如：账户创建时间、积分充值总额、积分消费总额等与写操作无关的信息，而这些展示信息
// 可以由 CQRS 的查询模型提供。
type ormCqrsAccount struct {
	baseAccount
	events *model.EventRecorder
}

// RegisterOrmCqrsAccount 注册一个 ORM & CQRS 风格的积分账户
func RegisterOrmCqrsAccount(ownerID model.StringID) Account {
	a := &ormCqrsAccount{events: model.NewEventRecorder(0)}
	a.id = model.IDGenerator.LongID()
	a.ownerID = ownerID
	a.events.RecordThan(OccurAccountCreated(a.id, ownerID))
	return a
}

func (a *ormCqrsAccount) Deposit(points Points) {
	a.points = a.points.Inc(points)
	a.events.RecordThan(occurDeposited(a.id, points))
}

func (a *ormCqrsAccount) Consume(points Points) error {
	if !a.points.GreaterThan(points) {
		return fmt.Errorf("当前账户积分为 %d 不足消费额 %d", a.points, points)
	}
	a.points = a.points.Inc(points)
	a.events.RecordThan(occurConsumed(a.id, points))
	return nil
}

func (a ormCqrsAccount) Version() uint {
	return a.events.LastVersion()
}

func (a *ormCqrsAccount) CommitEvents(publishers ...model.EventPublisher) {
	a.events.CommitToPublisher(publishers...)
}
