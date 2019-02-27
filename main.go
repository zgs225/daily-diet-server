package main

import (
	"log"
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
}

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err)
	}
}
