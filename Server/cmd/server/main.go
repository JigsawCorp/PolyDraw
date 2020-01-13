package main

import (
	"log"

	"gitlab.com/jigsawcorp/log3900/internal/rest"
	"gitlab.com/jigsawcorp/log3900/pkg/graceful"
)

func main() {
	restServer := &rest.Server{}

	graceful.Register(restServer.Shutdown, "REST server")
	graceful.ListenSIG()

	hRestServer := make(chan bool)
	go func() {
		restServer.Initialize()
		restServer.Run(":3000")
		hRestServer <- true
	}()

	log.Printf("Server is starting jobs!")

	//TODO Launch other servers and handles

	<-hRestServer
}
