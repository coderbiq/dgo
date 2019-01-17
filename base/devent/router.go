package devent

type simpleRouter struct {
	consumers map[string][]Consumer
}

// SimpleRouter 创建一个简单的事件路由
func SimpleRouter(routes map[string][]Consumer) Router {
	router := &simpleRouter{consumers: routes}
	return router
}

func (router *simpleRouter) Consumers(eventName string) ([]Consumer, bool) {
	consumers, has := router.consumers[eventName]
	return consumers, has
}
