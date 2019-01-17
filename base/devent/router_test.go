package devent_test

import (
	"testing"

	"github.com/coderbiq/dgo/base/devent"
)

func TestRegexRouter(t *testing.T) {
	routes := map[string][]devent.Consumer{
		"foo.*": []devent.Consumer{
			devent.ConsumerFunc(func(e devent.Event) {}),
		},
		"foo.name": []devent.Consumer{
			devent.ConsumerFunc(func(e devent.Event) {}),
			devent.ConsumerFunc(func(e devent.Event) {}),
		},
		"bar.*": []devent.Consumer{
			devent.ConsumerFunc(func(e devent.Event) {}),
		},
	}

	router := devent.RegexRouter(routes)
	if consumers, has := router.Consumers("foo.name"); !has ||
		len(consumers) != 3 {
		t.FailNow()
	}
	if consumers, has := router.Consumers("foo.other"); !has ||
		len(consumers) != 1 {
		t.FailNow()
	}
}
