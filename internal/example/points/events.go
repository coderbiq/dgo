package points

import "github.com/coderbiq/dgo/model"

const (
	// AccountCreated 积分账户创建事件
	AccountCreated = "accountCreated"
	// AccountDeposited 积分账户充值事件
	AccountDeposited = "accountDeposited"
	// AccountConsumed 积分账户消费事件
	AccountConsumed = "accountConsumed"
)

type (
	// AccountCreatedPayload 积分账户创建事件信息
	AccountCreatedPayload struct {
		ownerID CustomerID
	}
	// AccountDepositedPayload 积分账户充值事件信息
	AccountDepositedPayload struct {
		points Points
	}
	// AccountConsumedPayload 积分账户消费事件信息
	AccountConsumedPayload struct {
		points Points
	}
)

// NewAccountCreatedEvent 返回一个新的积分账户创建成功事件
func NewAccountCreatedEvent(aid AccountID, ownerID CustomerID) model.DomainEvent {
	return model.OccurEvent(
		aid,
		AccountCreated,
		&AccountCreatedPayload{ownerID: ownerID})
}

// OwnerID 返回新建积分账户所属的会员唯一标识
func (p AccountCreatedPayload) OwnerID() CustomerID {
	return p.ownerID
}

func newDepositedEvent(aid AccountID, points Points) model.DomainEvent {
	return model.OccurEvent(
		aid,
		AccountDeposited,
		&AccountDepositedPayload{points: points})
}

// Points 返回积分账户充值金额
func (p AccountDepositedPayload) Points() Points {
	return p.points
}

func newConsumedEvent(aid AccountID, points Points) model.DomainEvent {
	return model.OccurEvent(
		aid,
		AccountConsumed,
		&AccountConsumedPayload{points: points})
}

// Points 返回积分账户消息金额
func (p AccountConsumedPayload) Points() Points {
	return p.points
}
