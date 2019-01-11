package dgo_test

import (
	"testing"

	"github.com/coderbiq/dgo"
	"github.com/stretchr/testify/assert"
)

func TestImplementIdentity(t *testing.T) {
	assert.Implements(t, (*dgo.Identity)(nil), dgo.StringID("testId"))
	assert.Implements(t, (*dgo.Identity)(nil), dgo.LongID(1))
}

func TestToString(t *testing.T) {
	assert.Equal(t, "testId", dgo.StringID("testId").String())
	assert.Equal(t, "1", dgo.LongID(1).String())
}

func TestEquals(t *testing.T) {
	id1 := dgo.StringID("testId")
	id2 := dgo.StringID("testId")
	assert.True(t, id1.Equal(id2))
	assert.False(t, id2.Equal(dgo.StringID("other")))

	id3 := dgo.LongID(1)
	id4 := dgo.LongID(1)
	assert.True(t, id3.Equal(id4))
	assert.False(t, id4.Equal(dgo.LongID(2)))

	assert.False(t, id1.Equal(id3))
	assert.False(t, id4.Equal(id2))
}
