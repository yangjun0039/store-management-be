package network

import (
	"net/http"
)

type handlerBuilder struct {
	middlewares [] func(http.Handler) http.Handler
	adapter     func(Handler) http.Handler
	thens       []func(Handler) Handler
}

func NewHandlerBuilder(middlewares ...func(http.Handler) http.Handler) *handlerBuilder {
	builder := &handlerBuilder{middlewares: middlewares}
	return builder
}

func (self *handlerBuilder) Adapt(middleware func(Handler) http.Handler) *handlerBuilder {
	self.adapter = middleware
	return self
}

func (self *handlerBuilder) Then(handlers ...func(Handler) Handler) *handlerBuilder {
	self.thens = append(self.thens, handlers...)
	return self
}

func (self *handlerBuilder) Run(handler Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := len(self.middlewares)
		var thenHandler Handler
		var endHandler http.Handler
		if self.adapter != nil {
			count := len(self.thens)
			if count > 0 {
				thenHandler = self.thens[count-1](handler)
				for i := count - 2; i >= 0; i-- {
					thenHandler = self.thens[i](thenHandler)
				}
			} else {
				thenHandler = handler
			}
			endHandler = self.adapter(thenHandler)
		} else {
			panic("lack adapter")
		}

		if count == 0 {
			endHandler.ServeHTTP(w, r)
		} else {
			//endHandler = self.middlewares[count - 1](endHandler)
			for i := count - 1; i >= 0; i -- {
				endHandler = self.middlewares[i](endHandler)
			}
			endHandler.ServeHTTP(w, r)
		}
	})

}
