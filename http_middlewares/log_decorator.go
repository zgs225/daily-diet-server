package http_middlewares

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTPLogDecorator 记录每次 HTTP 请求装饰器
type HTTPLogDecorator struct {
	Logger *log.Entry
	next   http.Handler
}

func (d *HTTPLogDecorator) Decorate(next http.Handler) http.Handler {
	d.next = next
	if d.Logger == nil {
		return next
	}
	return d
}

func (d *HTTPLogDecorator) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	lw := &saveCodeResponseWriter{w, http.StatusOK}
	defer func(b time.Time) {
		d.Logger.WithFields(log.Fields{
			"Method":   request.Method,
			"URI":      request.URL.RequestURI(),
			"Code":     lw.code,
			"Duration": time.Since(b),
		}).Info(http.StatusText(lw.code))
	}(time.Now())
	d.next.ServeHTTP(lw, request)
}

// saveCodeResponseWriter

type saveCodeResponseWriter struct {
	http.ResponseWriter
	code int
}

func (w *saveCodeResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
