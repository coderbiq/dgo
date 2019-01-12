package eventsourcing_test

import (
	"testing"

	"github.com/coderbiq/dgo/internal/example"
	"github.com/coderbiq/dgo/model"
	"github.com/stretchr/testify/assert"
)

func TestSourcedEventRecorder(t *testing.T) {
	id := model.StringID("testId")
	text := "test text"
	aggregate := example.PostSourcedTodo(id, text)

	assert := assert.New(t)
	assert.Equal(id, aggregate.ID())
	assert.Equal(text, aggregate.Text())
	assert.Equal(1, int(aggregate.Version()))
	assert.Equal(1, len(aggregate.RecordedEvents()))
	assert.Equal(1, int(aggregate.RecordedEvents()[0].Version()))
}
