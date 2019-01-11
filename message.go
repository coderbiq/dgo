package dgo

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

// IdentityGenerator generate identity
var IdentityGenerator identityGenerator = defIdentityGenerator

type identityGenerator func() Identity

// Payload message payload
type Payload interface {
	JSON() ([]byte, error)
}

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

func defIdentityGenerator() Identity {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return LongID(node.Generate())
}
