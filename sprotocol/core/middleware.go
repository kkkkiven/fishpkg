package core

var middleware []func(Handler) Handler

func init() {
	middleware = make([]func(Handler) Handler, 0)
}

// Use 添加中间件
func Use(mws ...func(Handler) Handler) {
	middleware = append(middleware, mws...)
}

func chain(endpoint Handler) Handler {
	if len(middleware) == 0 {
		return endpoint
	}

	h := middleware[len(middleware)-1](endpoint)
	for i := len(middleware) - 2; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}
