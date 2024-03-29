package devent

// Recorder 定义领域事件记录器
//
// 事件记录器的目的是为了让聚合专注于领域事件的生产而解耦聚合和事件存储、事件发布机制的关系。
// 在每个聚合内部都可以持有一个事件记录器，当聚合内部产生领域事件后只需要将事件记录到记录器中，聚合
// 不需要关心被记录的事件持续会被如何处理。而在应用层中一个应用服务在调用聚合的命令方法执行完业务
// 操作后，可以将聚合内部暂存的已发生事件提交到事件存储进行持久化或提交到消息发布器进行发布。
type Recorder struct {
	version        uint
	recordedEvents []Event
}

// NewRecorder 创建一个事件记录器，version 参数为持有事件记录器的聚合当前的最新版本。
func NewRecorder(version uint) *Recorder {
	return &Recorder{
		version:        version,
		recordedEvents: []Event{},
	}
}

// RecordThan 记录一个已发生的领域事件
func (r *Recorder) RecordThan(event Event) {
	if aggregateChanged, ok := event.(aggregateChanged); ok {
		r.version++
		aggregateChanged.init()
		aggregateChanged.withVersion(r.version)
	}
	if err := ValidEvent(event); err != nil {
		panic(err)
	}
	r.recordedEvents = append(r.recordedEvents, event)
}

// CommitToPublisher 将记录器暂存的事件提交到事件发布器并清空暂存列表，持有事件记录器的聚合可
// 以将事件提交行为直接代理给当前方法。
func (r *Recorder) CommitToPublisher(publishers ...Publisher) {
	for _, publisher := range publishers {
		publisher.Publish(r.recordedEvents...)
	}
	r.recordedEvents = []Event{}
}

// RecordedEvents 返回当前暂存中记录的所有领域事件
func (r *Recorder) RecordedEvents() []Event {
	return r.recordedEvents
}

// LastVersion 返回根据事件记录计算的当前聚合最新版本
func (r *Recorder) LastVersion() uint {
	return r.version
}
