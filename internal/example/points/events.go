package points

import (
	"encoding/json"

	"github.com/coderbiq/dgo/model"
)

const (
	// AccountCreatedEvent 积分账户创建事件
	AccountCreatedEvent = "accountCreated"
	// AccountDepositedEvent 积分账户充值事件
	AccountDepositedEvent = "accountDeposited"
	// AccountConsumedEvent 积分账户消费事件
	AccountConsumedEvent = "accountConsumed"
)

type (
	// AccountCreated 积分账户创建事件信息
	AccountCreated interface {
		model.DomainEvent
		OwnerID() CustomerID
	}
	// AccountDeposited 积分账户充值事件信息
	AccountDeposited interface {
		model.DomainEvent
		Points() Points
	}
	// AccountConsumed 积分账户消费事件信息
	AccountConsumed interface {
		model.DomainEvent
		Points() Points
	}
)

type accountCreated struct {
	model.AggregateChanged
	// ownerID CustomerID
}

// OccurAccountCreated 返回一个新的积分账户创建成功事件
func OccurAccountCreated(aid AccountID, ownerID CustomerID) AccountCreated {
	e := accountCreated{
		AggregateChanged: model.AggregateChanged{
			Payload: map[string]interface{}{
				"aggregateId": aid.String(),
				"ownerId":     ownerID.String(),
			},
		},
	}
	e.Init()
	return e
	// return model.OccurEvent(
	// 	aid,
	// 	&accountCreated{ownerID: ownerID}).(AccountCreated)
}

func AccountCreatedFromJSON(data []byte) (AccountCreated, error) {
	e := &accountCreated{}
	if err := json.Unmarshal(data, e); err != nil {
		return nil, err
	}
	return e, nil
}

// Name 返回积分账户创建成功事件名称
func (p accountCreated) Name() string {
	return AccountCreatedEvent
}

// OwnerID 返回新建积分账户所属的会员唯一标识
func (p accountCreated) OwnerID() CustomerID {
	return model.IDFromInterface(p.Payload["ownerId"])
	// return p.ownerID
}

type accountDeposited struct {
	model.AggregateChanged
	points Points
}

func occurDeposited(aid AccountID, points Points) AccountDeposited {
	return model.OccurEvent(
		aid,
		&accountDeposited{points: points}).(AccountDeposited)
}

// Name 返回积分账户充值事件名称
func (p accountDeposited) Name() string {
	return AccountDepositedEvent
}

// Points 返回积分账户充值金额
func (p accountDeposited) Points() Points {
	return p.points
}

type accountConsumed struct {
	model.AggregateChanged
	points Points
}

func occurConsumed(aid AccountID, points Points) AccountConsumed {
	return model.OccurEvent(aid, &accountConsumed{points: points}).(AccountConsumed)
}

// Name 返回积分账户消费事件名称
func (p accountConsumed) Name() string {
	return AccountConsumedEvent
}

// Points 返回积分账户消息金额
func (p accountConsumed) Points() Points {
	return p.points
}
