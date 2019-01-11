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

func (suite *eventRecorderTestSuite) TestInAggregate() {
	example.PostTodo(dgo.StringID("testId"), "test text")
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
