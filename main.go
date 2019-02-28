package main

import (
	log "github.com/sirupsen/logrus"
)

var (
	app *HTTPApplication
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:             true,
		EnvironmentOverrideColors: true,
	})

	app = &HTTPApplication{
		Name: "daily-diet",
		Host: "0.0.0.0",
		Port: 10089,
	}
}

func main() {
	app.HandleRouter(api_route_func)

	if err := app.Run(); err != nil {
		log.Panic(err)
	}
}
