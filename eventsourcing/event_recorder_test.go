package eventsourcing_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSourcedEventRecorder(t *testing.T) {
	account := RegisterAccount("test account")

	assert := assert.New(t)
	assert.Equal("test account", account.Name)
	assert.NotEmpty(account.ID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	publisher := mocks.NewMockEventPublisher(ctrl)
	publisher.EXPECT().Publish(gomock.Any()).Times(1)
	account.CommitEvents(publisher)
}
