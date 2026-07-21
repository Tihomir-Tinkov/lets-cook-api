package routes

type RouterOptFunc func(r *Router)

func WithLoggingMiddleware() RouterOptFunc {
	return func(r *Router) {
		r.loggerMiddleware = true
	}
}

func WithPrefix(prefix string) RouterOptFunc {
	return func(r *Router) {
		r.prefix = prefix
	}
}
