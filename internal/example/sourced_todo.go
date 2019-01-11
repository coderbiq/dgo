package example

import (
	"github.com/coderbiq/dgo/eventsourcing"
	"github.com/coderbiq/dgo/model"
)

// SourcedTodo 实现了 EventSourcing 的 Todo 聚合根
type SourcedTodo struct {
	events *eventsourcing.EventRecorder
	id     TodoID
	text   string
}

// PostSourcedTodo 新建一条待办事项
func PostSourcedTodo(id TodoID, test string) *SourcedTodo {
	todo := &SourcedTodo{}
	todo.events = eventsourcing.EventRecorderFromSourced(todo, 0)
	todo.events.RecordThan(
		model.OccurDomainEvent(id, TodoCreated, NewTodoCreatedPayload(test)))
	return todo
}

// ID 返回当前聚合的唯一标识
func (t SourcedTodo) ID() model.Identity {
	return t.id
}

// Text 返回待办描述
func (t SourcedTodo) Text() string {
	return t.text
}

// Version 返回聚合版本
func (t SourcedTodo) Version() uint {
	return t.events.LastVersion()
}

// RecordedEvents 返回当前聚合暂存的所有领域事件
func (t SourcedTodo) RecordedEvents() []model.DomainEvent {
	return t.events.RecordedEvents()
}

// Apply 应用领域事件构建当前聚合状态
func (t *SourcedTodo) Apply(event model.DomainEvent) {
	switch event.Name() {
	case TodoCreated:
		payload := event.Payload().(TodoCreatedPayload)
		t.id = event.AggregateID()
		t.text = payload.Text()
		break
	}
}
