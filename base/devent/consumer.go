package devent

// ConsumerFunc 包装一个领域事件处理函数为领域事件消费者
func ConsumerFunc(handle func(Event)) Consumer {
	return &simpleConsumer{handle: handle}
}

type simpleConsumer struct {
	handle func(Event)
}

func (consume *simpleConsumer) Handle(event Event) {
	consume.handle(event)
}
