package model_test

import (
	"encoding/json"
	"testing"

	"github.com/coderbiq/dgo/internal/example/points"
	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	assert := assert.New(t)

	e := points.OccurAccountCreated(model.IdentityGenerator(), model.IdentityGenerator())
	data, err := json.Marshal(e)
	assert.Empty(err)

	e2, err := points.AccountCreatedFromJSON(data)
	assert.Empty(err)

	data2, err := json.Marshal(e2)
	assert.Equal(data, data2)
}
