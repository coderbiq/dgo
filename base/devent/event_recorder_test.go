package devent_test

import (
	"testing"

	"github.com/coderbiq/dgo/base/devent"
	"github.com/coderbiq/dgo/base/vo"
	"github.com/coderbiq/dgo/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type eventRecorderTestSuite struct {
	suite.Suite
	recorder *devent.EventRecorder
}

func (suite *eventRecorderTestSuite) SetupTest() {
	suite.recorder = devent.NewEventRecorder(0)
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
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	publisher := mocks.NewMockEventPublisher(ctrl)
	publisher.EXPECT().
		Publish(gomock.AssignableToTypeOf(&AccountCreated{})).
		Times(1)

	account := RegisterAccount("account name")
	account.CommitEvents(publisher)
}

func (suite *eventRecorderTestSuite) newEvent() devent.DomainEvent {
	return occurAccountCreate(
		vo.IDGenerator.LongID(),
		"test account")
}

func TestEventRecorderSuite(t *testing.T) {
	suite.Run(t, new(eventRecorderTestSuite))
}
