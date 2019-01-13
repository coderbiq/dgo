package model

import (
	"encoding/json"
	"errors"
	"time"
)

// DomainEvent 定义领域事件外观
type DomainEvent interface {
	ID() Identity
	Name() string
	AggregateID() Identity
	Version() uint
	CreatedAt() time.Time
}

// EventPublisher 定义事件发布器
//
// 事件发布器可以将领域事件发布到消息中间件或事件存储，具体如何发布由发布器内部进行实现。
type EventPublisher interface {
	// Publish 发布领域事件
	//
	// 发布器内部应该自行维护失败状态。例如：提交消息中间件失败后可以将需要发布的事件暂存，待检查
	// 到消息中间件恢复后再继续发布。
	Publish(events ...DomainEvent)
}

// EventProducer 定义领域事件生产者
type EventProducer interface {
	CommitEvents(publishers ...EventPublisher)
}

// ValidDomainEvent 验证一个领域事件是否完整
func ValidDomainEvent(event DomainEvent) error {
	if event.ID() == nil || event.ID().Empty() {
		return errors.New("事件标识不能为空")
	}
	if event.Name() == "" {
		return errors.New("事件名称不能为空")
	}
	if event.AggregateID() == nil || event.AggregateID().Empty() {
		return errors.New("发生领域事件的聚合标识不能为空")
	}
	if event.Version() < 1 {
		return errors.New("无效的领域事件版本")
	}
	if event.CreatedAt().IsZero() {
		return errors.New("领域事件发生时间不能为空")
	}
	return nil
}

// OccurEvent 组装一个已发生的领域事件
func OccurEvent(aggregateID Identity, event DomainEvent) DomainEvent {
	if changed, ok := event.(aggregateChanged); ok {
		changed.WithAggregateID(aggregateID)
		changed.Init()
	}
	return event
}

// AggregateChanged 提供对 DomainEvent 基本外观的实现
type AggregateChanged struct {
	Payload map[string]interface{}
}

type aggregateChanged interface {
	Init()
	WithVersion(version uint)
	WithAggregateID(id Identity)
}

// ID 返回聚合变更事件标识
func (ac AggregateChanged) ID() Identity {
	return IDFromInterface(ac.Payload["id"])
}

// AggregateID 返回发生变更的聚合标识
func (ac AggregateChanged) AggregateID() Identity {
	return IDFromInterface(ac.Payload["aggregateId"])
}

// Version 返回聚合变更时的版本号
func (ac AggregateChanged) Version() uint {
	if v, has := ac.Payload["version"]; has {
		return v.(uint)
	}
	return 0
}

// CreatedAt 返回聚合变更的时间
func (ac AggregateChanged) CreatedAt() time.Time {
	return ac.Payload["createdAt"].(time.Time)
}

// Init 初始化变更事件
func (ac *AggregateChanged) Init() {
	if ac.Payload == nil {
		ac.Payload = map[string]interface{}{}
	}
	if _, has := ac.Payload["id"]; !has {
		ac.Payload["id"] = IdentityGenerator().String()
	}
	if _, has := ac.Payload["createdAt"]; !has {
		ac.Payload["createdAt"] = time.Now()
	}
}

// WithAggregateID 指定发生变更的聚合标识
func (ac *AggregateChanged) WithAggregateID(id Identity) {
	ac.Payload["aggregateId"] = id.String()
}

// WithVersion 指定聚合变更版本
func (ac *AggregateChanged) WithVersion(version uint) {
	ac.Payload["version"] = version
}

// MarshalJSON json 序列化
func (ac AggregateChanged) MarshalJSON() ([]byte, error) {
	return json.Marshal(ac.Payload)
}

// UnmarshalJSON json 反序列化
func (ac *AggregateChanged) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	ac.Payload = payload
	return nil
}
