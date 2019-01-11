package example

import "github.com/coderbiq/dgo"

// OrmTodo aggregate root model
type OrmTodo struct {
	events *dgo.EventRecorder
	id     TodoID
	text   string
}

// PostOrmTodo post todo
func PostOrmTodo(id TodoID, text string) *OrmTodo {
	todo := &OrmTodo{events: dgo.NewEventRecorder(0)}
	todo.id = id
	todo.text = text
	todo.events.RecordThan(
		dgo.OccurDomainEvent(id, TodoCreated, NewTodoCreatedPayload(text)))
	return todo
}

// Id return aggregate id
func (t OrmTodo) Id() dgo.Identity {
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
func (t OrmTodo) RecordedEvents() []dgo.DomainEvent {
	return t.events.RecordedEvents()
}
