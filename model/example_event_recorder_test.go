package model_test

import "github.com/coderbiq/dgo/model"

// Account 定义一个持有事件记录器的聚合模型
type Account struct {
	// 在聚合内部持有一个事件记录器实例，用例记录聚合在本次会话过程中发生的领域事件
	events *model.EventRecorder

	ID   model.LongID
	Name string
}

// CommitEvents 持有事件记录器的聚合应该实现 EventProducer 接口，让应用服务可以将聚合内产生的
// 事件提交给一些事件发布器。事件发布器可以将事件提交到消息中间件、持久化存储、事件总线等进行后续操作
func (account *Account) CommitEvents(publishers ...model.EventPublisher) {
	account.events.CommitToPublisher(publishers...)
}

func RegisterAccount(name string) *Account {
	account := &Account{events: model.NewEventRecorder(0)}
	account.ID = model.IDGenerator.LongID()
	account.Name = name
	// 将账户注册成功事件记录到事件记录器
	account.events.RecordThan(occurAccountCreate(account.ID, account.Name))
	return account
}

// AccountCreated 定义聚合模型创建成功事件
type AccountCreated struct {
	// 通过组合 AggreateChanged 获取领域事件的基本能力
	model.AggregateChanged

	AccountID   model.LongID `json:"aggregateId"`
	AccountName string       `json:"accountName"`
}

func occurAccountCreate(aid model.LongID, name string) *AccountCreated {
	return model.OccurAggregateChanged(
		"accountCreated",
		&AccountCreated{
			AccountID:   aid,
			AccountName: name,
		}).(*AccountCreated)
}

func (event AccountCreated) AggregateID() model.Identity {
	return event.AccountID
}
