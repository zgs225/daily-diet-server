package http_middlewares

import (
	"net/http"
)

// HTTPHandlerDecorator http.Handler 装饰器
type HTTPHandlerDecorator interface {
	Decorate(http.Handler) http.Handler
}
