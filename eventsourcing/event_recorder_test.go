package eventsourcing_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/example/points"
	"github.com/coderbiq/dgo/internal/mocks"
	"github.com/coderbiq/dgo/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSourcedEventRecorder(t *testing.T) {
	ownerID := model.IDGenerator.StringID()
	account := points.RegisterSourcedAccount(ownerID)

	assert := assert.New(t)
	assert.True(ownerID.Equal(account.OwnerID()))
	assert.False(account.ID().Empty())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	publisher := mocks.NewMockEventPublisher(ctrl)
	publisher.EXPECT().Publish(gomock.Any()).Times(1)
	account.(model.EventProducer).CommitEvents(publisher)
}
