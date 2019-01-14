package model_test

import (
	"encoding/json"
	"testing"

	"github.com/coderbiq/dgo/internal/example/points"
	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/assert"
)

func TestDomainEventToJson(t *testing.T) {
	assert := assert.New(t)

	e := points.OccurAccountCreated(
		model.IDGenerator.LongID(),
		model.IDGenerator.StringID())

	data, err := json.Marshal(e)
	assert.Empty(err)

	m := make(map[string]interface{})
	assert.Empty(json.Unmarshal(data, &m))
	assert.Contains(m, "id")
	assert.Contains(m, "name")
	assert.Contains(m, "aggregateId")
	assert.Contains(m, "version")
	assert.Contains(m, "createdAt")
	assert.Contains(m, "ownerId")

	e2, err := points.AccountCreatedFromJSON(data)
	assert.Empty(err)

	data2, _ := json.Marshal(e2)
	assert.Equal(string(data), string(data2))
}
