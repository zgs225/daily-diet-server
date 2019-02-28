package main

import (
	"fmt"
	"io"
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

	server  *http.Server
	logger  *log.Entry
	handler http.Handler
	router  *http.ServeMux
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

func (app *HTTPApplication) init() error {
	app.initLogger()
	app.initHTTPHandler()

	if err := app.initServer(); err != nil {
		return err
	}

	app.router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, app.Name)
	})

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
