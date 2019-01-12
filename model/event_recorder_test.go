package model_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/example"
	"github.com/coderbiq/dgo/internal/mocks"
	"github.com/coderbiq/dgo/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type eventRecorderTestSuite struct {
	suite.Suite
	recorder *model.EventRecorder
}

func (suite *eventRecorderTestSuite) SetupTest() {
	suite.recorder = model.NewEventRecorder(0)
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

func (suite *eventRecorderTestSuite) TestCommitToPublisher() {
	suite.recorder.RecordThan(suite.newEvent())
	suite.recorder.RecordThan(suite.newEvent())
	events := suite.recorder.RecordedEvents()

	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	publisher1 := mocks.NewMockEventPublisher(ctrl)
	publisher1.EXPECT().Publish(events[0], events[1]).Times(1)
	publisher2 := mocks.NewMockEventPublisher(ctrl)
	publisher2.EXPECT().Publish(events[0], events[1]).Times(1)

	suite.recorder.CommitToPublisher(publisher1, publisher2)
}

func (suite *eventRecorderTestSuite) TestInOrmAggregate() {
	id := model.StringID("testId")
	text := "test text"
	aggregate := example.PostOrmTodo(id, text)
	suite.Equal(id, aggregate.ID())
	suite.Equal(text, aggregate.Text())
	suite.Equal(1, int(aggregate.Version()))
	suite.Equal(1, len(aggregate.RecordedEvents()))
	suite.Equal(1, int(aggregate.RecordedEvents()[0].Version()))
}

func (suite *eventRecorderTestSuite) newEvent() model.DomainEvent {
	return model.OccurDomainEvent(
		model.StringID("testAggregateId"),
		example.TodoCreated,
		example.NewTodoCreatedPayload("test text"))
}

func TestEventRecorderSuite(t *testing.T) {
	suite.Run(t, new(eventRecorderTestSuite))
}
