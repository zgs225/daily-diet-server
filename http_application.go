package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// HTTPApplication 用于启动 HTTP 服务
type HTTPApplication struct {
	Name string
	Host string
	Port int

	server *http.Server
}

func (app *HTTPApplication) Run() error {
	c := make(chan error)
	go func() {
		if err := app.init(); err != nil {
			c <- err
			return
		}
		addr := fmt.Sprintf("%s:%d", app.Host, app.Port)
		log.Printf("%s: http server starting at %s", app.Name, addr)
		c <- http.ListenAndServe(addr, http.DefaultServeMux)
	}()
	return <-c
}

func (app *HTTPApplication) init() error {
	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, app.Name)
	})

	return nil
}
