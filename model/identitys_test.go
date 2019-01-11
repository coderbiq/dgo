package model_test

import (
	"testing"

	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/assert"
)

func TestImplementIdentity(t *testing.T) {
	assert.Implements(t, (*model.Identity)(nil), model.StringID("testId"))
	assert.Implements(t, (*model.Identity)(nil), model.LongID(1))
}

func TestToString(t *testing.T) {
	assert.Equal(t, "testId", model.StringID("testId").String())
	assert.Equal(t, "1", model.LongID(1).String())
}

func TestEquals(t *testing.T) {
	id1 := model.StringID("testId")
	id2 := model.StringID("testId")
	assert.True(t, id1.Equal(id2))
	assert.False(t, id2.Equal(model.StringID("other")))

	id3 := model.LongID(1)
	id4 := model.LongID(1)
	assert.True(t, id3.Equal(id4))
	assert.False(t, id4.Equal(model.LongID(2)))

	assert.False(t, id1.Equal(id3))
	assert.False(t, id4.Equal(id2))
}
