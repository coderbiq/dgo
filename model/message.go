package model

import (
	"time"
)

// Message 定义消息格式
type Message struct {
	id      Identity
	name    string
	payload interface{}
	created time.Time
}

// NewMessage 创建一个新的消息
func NewMessage(name string, payload interface{}) Message {
	return Message{
		id:      IdentityGenerator(),
		name:    name,
		payload: payload,
		created: time.Now(),
	}
}

// ID 返回消息唯一标识
func (e Message) ID() Identity {
	return e.id
}

// Name 返回消息名称
func (e Message) Name() string {
	return e.name
}

// Payload 返回消息有效载荷
func (e Message) Payload() interface{} {
	return e.payload
}

// CreatedAt 返回消息的创建时间
func (e Message) CreatedAt() time.Time {
	return e.created
}
