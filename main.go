package main

import (
	log "github.com/sirupsen/logrus"
)

var (
	app Application
)

func init() {
	app = &HTTPApplication{
		Name: "daily-diet",
		Host: "0.0.0.0",
		Port: 10089,
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:             true,
		EnvironmentOverrideColors: true,
	})
}

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err)
	}
}
