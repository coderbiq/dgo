package base

import (
	"time"

	"github.com/coderbiq/dgo/internal/example/common"
	"github.com/coderbiq/dgo/model"
)

// Account 定义不使用 CQRS 的积分账户外观
type Account interface {
	common.Account

	DepositedPoints() common.Points
	ConsumedPoints() common.Points
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// AccountReadModel 积分展示读模型的实现
type AccountReadModel struct {
	DepPoints common.Points `json:"depositedPoints"`
	ConPoints common.Points `json:"consumedPoints"`
	Created   time.Time     `json:"createdAt"`
	Updated   time.Time     `json:"updatedAt"`
}

// DepositedPoints 返回积分账户历史充值总额
func (arm AccountReadModel) DepositedPoints() common.Points {
	return arm.DepPoints
}

// ConsumedPoints 返回积分账户历史消息总额
func (arm AccountReadModel) ConsumedPoints() common.Points {
	return arm.ConPoints
}

// CreatedAt 返回积分账户创建时间
func (arm AccountReadModel) CreatedAt() time.Time {
	return arm.Created
}

// UpdatedAt 返回积分账户最后变更时间
func (arm AccountReadModel) UpdatedAt() time.Time {
	return arm.Updated
}

type account struct {
	common.BaseAccount
	AccountReadModel
	events *model.EventRecorder
}

// RegisterAccount 为指定会员标识的会员注册一个新的积分账户
func RegisterAccount(ownerID model.StringID) Account {
	a := &account{events: model.NewEventRecorder(0)}
	a.Identity = model.IDGenerator.LongID()
	a.OwnerIdentity = ownerID
	a.DepPoints = common.Points(0)
	a.ConPoints = common.Points(0)
	a.Created = time.Now()
	a.Updated = time.Now()
	a.events.RecordThan(common.OccurAccountCreated(a.Identity, ownerID))
	return a
}

func (a account) Deposit(points common.Points) {
	a.CurPoints = a.CurPoints.Inc(points)
	a.DepPoints = a.DepPoints.Inc(points)
	a.Updated = time.Now()
	a.events.RecordThan(common.OccurDeposited(a.Identity, points))
}

func (a account) Consume(points common.Points) error {
	if !a.CurPoints.GreaterThan(points) {
		return common.PointsInsufficientError{
			Current: a.CurPoints,
			Expect:  points,
		}
	}
	a.CurPoints = a.CurPoints.Dec(points)
	a.ConPoints = a.ConPoints.Inc(points)
	a.Updated = time.Now()
	a.events.RecordThan(common.OccurConsumed(a.Identity, points))
	return nil
}
