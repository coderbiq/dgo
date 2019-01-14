package eventsourcing_test

import (
	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

// Account 定义一个持有事件记录器的聚合模型
type Account struct {
	// 聚合内部持有的是一个带有事件回放功能的事件记录器
	events *eventsourcing.EventRecorder

	ID   model.LongID
	Name string
}

// CommitEvents 持有事件记录器的聚合应该实现 EventProducer 接口，让应用服务可以将聚合内产生的
// 事件提交给一些事件发布器。事件发布器可以将事件提交到消息中间件、持久化存储、事件总线等进行后续操作
func (account *Account) CommitEvents(publishers ...model.EventPublisher) {
	account.events.CommitToPublisher(publishers...)
}

// 使用 EventSourcing 的聚合需要实现 EventSourcd 接口，
// Repository 和 EventRecorder 可以利用这个接口将领域事件应用到聚合以重建出聚合的状态
func (account *Account) Apply(event model.DomainEvent) {
	switch e := event.(type) {
	case *AccountCreated:
		account.ID = e.AccountID
		account.Name = e.AccountName
		break
	}
}

// EventSourcing 的聚合在业务命令中不用立即修改聚合的状态，而是可以在 Apply 方法中通过应用
// 事件的状态来构建聚合状态
func RegisterAccount(name string) *Account {
	account := new(Account)
	account.events = eventsourcing.EventRecorderFromSourced(account, 0)
	account.events.RecordThan(occurAccountCreate(
		model.IDGenerator.LongID(),
		name,
	))
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
