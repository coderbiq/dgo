package example

import "github.com/coderbiq/dgo"

// TodoID todo aggregate identity model
type TodoID dgo.Identity

// Todo aggregate root model
type Todo struct {
	ID   TodoID
	Text string
}
