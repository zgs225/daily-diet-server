package main

import (
	"io"
	"net/http"
)

func api_route_func(r *http.ServeMux) {
	r.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "hello world")
	})
}
