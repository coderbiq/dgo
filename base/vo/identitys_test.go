package vo_test

import (
	"testing"

	"github.com/coderbiq/dgo/base/vo"
	"github.com/stretchr/testify/assert"
)

func TestImplementIdentity(t *testing.T) {
	assert.Implements(t, (*vo.Identity)(nil), vo.StringID("testId"))
	assert.Implements(t, (*vo.Identity)(nil), vo.LongID(1))
}

func TestToString(t *testing.T) {
	assert.Equal(t, "testId", vo.StringID("testId").String())
	assert.Equal(t, "1", vo.LongID(1).String())
}

func TestEquals(t *testing.T) {
	id1 := vo.StringID("testId")
	id2 := vo.StringID("testId")
	assert.True(t, id1.Equal(id2))
	assert.False(t, id2.Equal(vo.StringID("other")))

	id3 := vo.LongID(1)
	id4 := vo.LongID(1)
	assert.True(t, id3.Equal(id4))
	assert.False(t, id4.Equal(vo.LongID(2)))

	assert.False(t, id1.Equal(id3))
	assert.False(t, id4.Equal(id2))
}
