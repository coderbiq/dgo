package model_test

import "github.com/coderbiq/dgo/model"

// Account 定义聚合模型，在聚合模型内部使用事件记录器用于存储聚合内部产生的各种领域事件
type Account struct {
	events *model.EventRecorder
	ID     model.LongID
	Name   string
}

// CommitEvents 方法让聚合外部可以将聚合内发生的事件提交到事件发布器
func (account *Account) CommitEvents(publishers ...model.EventPublisher) {
	account.events.CommitToPublisher(publishers...)
}

// 注册账户业务命令方法
//
// 在修改完聚合的内部状态后将产生的领域事件记录到事件记录器
func RegisterAccount(name string) *Account {
	account := &Account{events: model.NewEventRecorder(0)}
	account.ID = model.IDGenerator.LongID()
	account.Name = name
	// 将账户注册成功事件记录到事件记录器
	account.events.RecordThan(occurAccountCreate(account.ID, account.Name))
	return account
}

// AccountCreated 账户创建成功事件，存储领域事件的相关信息
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
	var eventBus model.EventPublisher
	account := RegisterAccount("test account")
	// 应用层中调用聚合执行完业务指令后，将聚合内部产生的领域事件发布到系统
	account.CommitEvents(eventBus)
}
