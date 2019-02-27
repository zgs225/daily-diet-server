package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTPApplication 用于启动 HTTP 服务
type HTTPApplication struct {
	Name string
	Host string
	Port int

	server *http.Server
	logger *log.Entry
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

	if err := app.initServer(); err != nil {
		return err
	}

	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
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

func (app *HTTPApplication) initServer() error {
	app.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.Host, app.Port),
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ConnState: func(conn net.Conn, stat http.ConnState) {
			app.logger.Debugf("A connection state changed: %s", stat.String())
		},
	}
	return nil
}

func (app *HTTPApplication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, req)
}
