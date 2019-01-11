package model

import (
	"time"
)

type baseMessage struct {
	id      Identity
	name    string
	payload Payload
	created time.Time
}

func newBaseMessage(name string, payload Payload) baseMessage {
	return baseMessage{
		id:      IdentityGenerator(),
		name:    name,
		payload: payload,
		created: time.Now(),
	}
}
