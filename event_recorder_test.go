package dgo_test

import (
	"testing"

	"github.com/coderbiq/dgo"
	"github.com/coderbiq/dgo/example"
	"github.com/stretchr/testify/suite"
)

type eventRecorderTestSuite struct {
	suite.Suite
	recorder *dgo.EventRecorder
}

func (suite *eventRecorderTestSuite) SetupTest() {
	suite.recorder = dgo.NewEventRecorder(0)
}

func (suite *eventRecorderTestSuite) TestDefaultStatus() {
	suite.Equal(0, int(suite.recorder.LastVersion()))
	suite.Empty(suite.recorder.RecordedEvents())
}

func (suite *eventRecorderTestSuite) TestRecord() {
	suite.recorder.RecordThan(suite.newEvent())
	suite.Equal(1, int(suite.recorder.LastVersion()))
	suite.Equal(1, len(suite.recorder.RecordedEvents()))
	suite.Equal(1, int(suite.recorder.RecordedEvents()[0].Version()))

	suite.recorder.RecordThan(suite.newEvent())
	suite.Equal(2, int(suite.recorder.LastVersion()))
	suite.Equal(2, len(suite.recorder.RecordedEvents()))
	suite.Equal(1, int(suite.recorder.RecordedEvents()[0].Version()))
	suite.Equal(2, int(suite.recorder.RecordedEvents()[1].Version()))

}

func (suite *eventRecorderTestSuite) TestInOrmAggregate() {
	id := dgo.StringID("testId")
	text := "test text"
	aggregate := example.PostOrmTodo(id, text)
	suite.Equal(id, aggregate.ID())
	suite.Equal(text, aggregate.Text())
	suite.Equal(1, int(aggregate.Version()))
	suite.Equal(1, len(aggregate.RecordedEvents()))
	suite.Equal(1, int(aggregate.RecordedEvents()[0].Version()))
}

func (suite *eventRecorderTestSuite) TestInSourcedAggregate() {
	id := dgo.StringID("testId")
	text := "test text"
	aggregate := example.PostSourcedTodo(id, text)
	suite.Equal(id, aggregate.ID())
	suite.Equal(text, aggregate.Text())
	suite.Equal(1, int(aggregate.Version()))
	suite.Equal(1, len(aggregate.RecordedEvents()))
	suite.Equal(1, int(aggregate.RecordedEvents()[0].Version()))
}

func (suite *eventRecorderTestSuite) newEvent() dgo.DomainEvent {
	return dgo.OccurDomainEvent(
		dgo.StringID("testAggregateId"),
		example.TodoCreated,
		example.NewTodoCreatedPayload("test text"))
}

func TestEventRecorderSuite(t *testing.T) {
	suite.Run(t, new(eventRecorderTestSuite))
}
