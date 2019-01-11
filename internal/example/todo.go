package example

import (
	"encoding/json"

	"github.com/coderbiq/dgo/model"
)

const (
	// TodoCreated event name
	TodoCreated = "todo.created"
)

// TodoID todo aggregate identity model
type TodoID model.Identity

// TodoCreatedPayload todo created domain event model
type TodoCreatedPayload struct {
	text string
}

// NewTodoCreatedPayload return new TodoCreatedPayload
func NewTodoCreatedPayload(text string) TodoCreatedPayload {
	return TodoCreatedPayload{text: text}
}

// Text return text of created todo
func (e TodoCreatedPayload) Text() string {
	return e.text
}

// JSON to json string
func (e TodoCreatedPayload) JSON() ([]byte, error) {
	return json.Marshal(e)
}
