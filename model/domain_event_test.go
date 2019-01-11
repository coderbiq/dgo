package model_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/example"
	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/suite"
)

type domainEventTestSuite struct {
	suite.Suite

	aid     model.Identity
	payload example.TodoCreatedPayload
}

func (suite *domainEventTestSuite) SetupTest() {
	suite.aid = model.StringID("testId")
	suite.payload = example.NewTodoCreatedPayload("test text")
}

func (suite *domainEventTestSuite) TestCreateEvent() {
	e := suite.newEvent()

	assert := suite.Assert()
	assert.Equal(suite.aid, e.AggregateID())
	assert.Equal(example.TodoCreated, e.Name())
	assert.Equal(suite.payload, e.Payload())
	assert.NotEmpty(e.CreatedAt())
	assert.Equal(0, int(e.Version()))
	assert.False(e.ID().Empty())

	created := e.Payload().(example.TodoCreatedPayload)
	assert.Equal("test text", created.Text())
}

func (suite *domainEventTestSuite) TestWithVersion() {
	e := suite.newEvent()
	e2 := e.WithVersin(2)
	suite.Equal(2, int(e2.Version()))
	suite.Equal(e.ID(), e2.ID())
	suite.Equal(e.AggregateID(), e2.AggregateID())
	suite.Equal(e.Name(), e2.Name())
	suite.Equal(e.CreatedAt(), e2.CreatedAt())
	suite.Equal(e.Payload(), e2.Payload())
}

func (suite *domainEventTestSuite) newEvent() model.DomainEvent {
	return model.OccurDomainEvent(
		suite.aid,
		example.TodoCreated,
		suite.payload)
}

func TestDomainEventSuite(t *testing.T) {
	suite.Run(t, new(domainEventTestSuite))
}
