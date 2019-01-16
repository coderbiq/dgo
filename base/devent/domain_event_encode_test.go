package devent_test

import (
	"encoding/json"
	"testing"

	"github.com/coderbiq/dgo/base/vo"
	"github.com/stretchr/testify/assert"
)

func TestDomainEventToJson(t *testing.T) {
	assert := assert.New(t)

	e := occurAccountCreate(vo.IDGenerator.LongID(), "test account")

	data, err := json.Marshal(e)
	assert.Empty(err)

	m := make(map[string]interface{})
	assert.Empty(json.Unmarshal(data, &m))
	assert.Contains(m, "id")
	assert.Contains(m, "name")
	assert.Contains(m, "aggregateId")
	assert.Contains(m, "version")
	assert.Contains(m, "createdAt")
	assert.Contains(m, "accountName")

	e2 := new(AccountCreated)
	assert.Empty(json.Unmarshal(data, e2))

	data2, _ := json.Marshal(e2)
	assert.Equal(string(data), string(data2))
}
