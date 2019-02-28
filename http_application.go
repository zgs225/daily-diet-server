package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zgs225/daily-diet-server/http_middlewares"
)

// HTTPApplication 用于启动 HTTP 服务
type HTTPApplication struct {
	Name string
	Host string
	Port int

	server      *http.Server
	logger      *log.Entry
	handler     http.Handler
	router      *http.ServeMux
	routerFuncs []func(*http.ServeMux)
}

func (app *HTTPApplication) Run() error {
	c := make(chan error)
	go func() {
		if err := app.init(); err != nil {
			c <- err
			return
		}
		app.logger.WithFields(log.Fields{
			"Host": app.Host,
			"Port": app.Port,
		}).Info("HTTP server starting")
		c <- app.server.ListenAndServe()
	}()
	return <-c
}

func (app *HTTPApplication) HandleRouteFunc(f func(*http.ServeMux)) {
	if app.routerFuncs == nil {
		app.routerFuncs = make([]func(*http.ServeMux), 0, 0)
	}
	app.routerFuncs = append(app.routerFuncs, f)
}

func (app *HTTPApplication) init() error {
	app.initLogger()
	app.initHTTPHandler()
	app.initServer()
	app.invokeRouterFuncs()
	return nil
}

func (app *HTTPApplication) initLogger() error {
	app.logger = log.WithFields(log.Fields{
		"app":  app.Name,
		"type": "HTTPApplication",
	})
	return nil
}

func (app *HTTPApplication) initHTTPHandler() error {
	app.router = http.NewServeMux()
	d := &http_middlewares.HTTPLogDecorator{Logger: app.logger}
	app.handler = d.Decorate(app.router)
	return nil
}

func (app *HTTPApplication) initServer() error {
	app.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.Host, app.Port),
		Handler:      app.handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ConnState: func(conn net.Conn, stat http.ConnState) {
			app.logger.Debugf("A connection state changed: %s", stat.String())
		},
	}
	return nil
}

func (app *HTTPApplication) invokeRouterFuncs() error {
	if app.routerFuncs == nil {
		return nil
	}
	for _, f := range app.routerFuncs {
		f(app.router)
	}
	return nil
}
