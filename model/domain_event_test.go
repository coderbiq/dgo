package model_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/example/points"
	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/suite"
)

type domainEventTestSuite struct {
	suite.Suite

	aid     model.Identity
	ownerId model.Identity
}

func (suite *domainEventTestSuite) SetupTest() {
	suite.aid = model.IdentityGenerator()
	suite.ownerId = model.IdentityGenerator()
}

func (suite *domainEventTestSuite) TestCreateEvent() {
	e := points.NewAccountCreatedEvent(suite.aid, suite.ownerId)

	assert := suite.Assert()
	assert.Equal(suite.aid, e.AggregateID())
	assert.Equal(points.AccountCreated, e.Name())
	assert.NotEmpty(e.CreatedAt())
	assert.Equal(0, int(e.Version()))
	assert.False(e.ID().Empty())

	suite.Equal(suite.ownerId, e.Payload().(*points.AccountCreatedPayload).OwnerID())
}

func (suite *domainEventTestSuite) TestWithVersion() {
	e := points.NewAccountCreatedEvent(suite.aid, suite.ownerId)

	e2 := e.WithVersion(2)
	suite.Equal(2, int(e2.Version()))
	suite.Equal(e.ID(), e2.ID())
	suite.Equal(e.AggregateID(), e2.AggregateID())
	suite.Equal(e.Name(), e2.Name())
	suite.Equal(e.CreatedAt(), e2.CreatedAt())
	suite.Equal(e.Payload(), e2.Payload())
}

func TestDomainEventSuite(t *testing.T) {
	suite.Run(t, new(domainEventTestSuite))
}
