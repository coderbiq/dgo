package devent

import "regexp"

type simpleRouter struct {
	consumers map[string][]Consumer
}

// SimpleRouter 创建一个简单的事件路由
func SimpleRouter(routes map[string][]Consumer) Router {
	return &simpleRouter{consumers: routes}
}

func (router *simpleRouter) Consumers(eventName string) ([]Consumer, bool) {
	consumers, has := router.consumers[eventName]
	return consumers, has
}

type regexRouter struct {
	consumers map[string][]Consumer
}

// RegexRouter 返回一个使用正则匹配的事件路由
func RegexRouter(routes map[string][]Consumer) Router {
	return &regexRouter{consumers: routes}
}

func (router *regexRouter) Consumers(eventName string) ([]Consumer, bool) {
	consumers := []Consumer{}
	has := false
	for pattern, cs := range router.consumers {
		if matched, err := regexp.MatchString(pattern, eventName); err == nil && matched {
			consumers = append(consumers, cs...)
			has = true
		}
	}
	return consumers, has
}
