package devent

import (
	"errors"
	"time"

	"github.com/coderbiq/dgo/base/vo"
)

type (
	// Event 定义领域事件外观
	Event interface {
		ID() vo.Identity
		Name() string
		AggregateID() vo.Identity
		Version() uint
		CreatedAt() time.Time
	}

	// Publisher 定义事件发布器
	//
	// 事件发布器可以将领域事件发布到消息中间件或事件存储，具体如何发布由发布器内部进行实现。
	Publisher interface {
		// Publish 发布领域事件
		//
		// 发布器内部应该自行维护失败状态。例如：提交消息中间件失败后可以将需要发布的事件暂存，待检查
		// 到消息中间件恢复后再继续发布。
		Publish(events ...Event)
	}

	// Producer 定义领域事件生产者
	Producer interface {
		CommitEvents(publishers ...Publisher)
	}

	// Consumer 定义领域事件消费者
	Consumer interface {
		Handle(Event)
	}

	// Router 根据事件名称返回订阅了该事件的所有消费者
	Router interface {
		Consumers(eventName string) ([]Consumer, bool)
	}

	// Bus 定义消息总线
	Bus interface {
		Publisher
		AddRouter(Router)
	}
)

// ValidEvent 验证一个领域事件是否完整
func ValidEvent(event Event) error {
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
