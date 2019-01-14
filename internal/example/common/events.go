package common

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
		OwnerID() model.StringID
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
	AccountID     model.LongID   `json:"aggregateId"`
	OwnerIdentity model.StringID `json:"ownerId"`
}

// OccurAccountCreated 返回一个新的积分账户创建成功事件
func OccurAccountCreated(aid model.LongID, ownerID model.StringID) AccountCreated {
	return model.OccurAggregateChanged(AccountCreatedEvent, &accountCreated{
		AccountID:     aid,
		OwnerIdentity: ownerID}).(AccountCreated)
}

// AccountCreatedFromJSON 通过 json 数据重建积分账户创建成功事件
func AccountCreatedFromJSON(data []byte) (AccountCreated, error) {
	e := &accountCreated{}
	if err := json.Unmarshal(data, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (p accountCreated) AggregateID() model.Identity {
	return p.AccountID
}

func (p accountCreated) OwnerID() model.StringID {
	return p.OwnerIdentity
}

type accountDeposited struct {
	model.AggregateChanged
	AccountID model.LongID
	points    Points
}

// OccurDeposited 返回一个积分账户充值成功事件
func OccurDeposited(aid model.LongID, points Points) AccountDeposited {
	return model.OccurAggregateChanged(
		AccountDepositedEvent,
		&accountDeposited{AccountID: aid, points: points}).(AccountDeposited)
}

func (p accountDeposited) AggregateID() model.Identity {
	return p.AccountID
}

func (p accountDeposited) Points() Points {
	return p.points
}

type accountConsumed struct {
	model.AggregateChanged
	AccountID model.LongID
	points    Points
}

// OccurConsumed 返回一个积分账户消费成功事件
func OccurConsumed(aid model.LongID, points Points) AccountConsumed {
	return model.OccurAggregateChanged(
		AccountConsumedEvent,
		&accountConsumed{AccountID: aid, points: points}).(AccountConsumed)
}

func (p accountConsumed) AggregateID() model.Identity {
	return p.AccountID
}

func (p accountConsumed) Points() Points {
	return p.points
}
