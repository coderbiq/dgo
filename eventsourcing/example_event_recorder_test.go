package eventsourcing_test

import (
	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

// Account 在聚合模型内部使用一个 EventSourcing 功能的事件记录器
type Account struct {
	events *eventsourcing.EventRecorder
	ID     model.LongID
	Name   string
}

// CommitEvents 方法可以用于将聚合内部产生的领域事件发布到事件存储
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

// AccountCreated 账户创建成功事件
type AccountCreated struct {
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

func ExampleEventRecorder() {
	var eventStore model.EventPublisher
	account := RegisterAccount("test account")
	// 将聚合内部产生的领域事件提交到事件存储中
	account.CommitEvents(eventStore)
}
