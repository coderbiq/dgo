package example

import (
	"github.com/coderbiq/dgo/model"
)

// OrmTodo aggregate root model
type OrmTodo struct {
	events *model.EventRecorder
	id     TodoID
	text   string
}

// PostOrmTodo post todo
func PostOrmTodo(id TodoID, text string) *OrmTodo {
	todo := &OrmTodo{events: model.NewEventRecorder(0)}
	todo.id = id
	todo.text = text
	todo.events.RecordThan(
		model.OccurDomainEvent(id, TodoCreated, NewTodoCreatedPayload(text)))
	return todo
}

// ID return aggregate id
func (t OrmTodo) ID() model.Identity {
	return t.id
}

// Text return todo text
func (t OrmTodo) Text() string {
	return t.text
}

// Version return aggregate version
func (t OrmTodo) Version() uint {
	return t.events.LastVersion()
}

// RecordedEvents return recorded domain events
func (t OrmTodo) RecordedEvents() []model.DomainEvent {
	return t.events.RecordedEvents()
}
