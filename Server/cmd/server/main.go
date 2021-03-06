package main

import (
	"fmt"
	"log"

	"gitlab.com/jigsawcorp/log3900/internal/language"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"

	"gitlab.com/jigsawcorp/log3900/internal/services/drawing"
	"gitlab.com/jigsawcorp/log3900/internal/services/lobby"
	"gitlab.com/jigsawcorp/log3900/internal/services/match"
	"gitlab.com/jigsawcorp/log3900/internal/services/potrace"
	redisservice "gitlab.com/jigsawcorp/log3900/internal/services/redis"
	"gitlab.com/jigsawcorp/log3900/internal/services/stats"
	"gitlab.com/jigsawcorp/log3900/internal/services/virtualplayer"
	"gitlab.com/jigsawcorp/log3900/pkg/geometry"

	"github.com/spf13/viper"
	"gitlab.com/jigsawcorp/log3900/internal/config"
	"gitlab.com/jigsawcorp/log3900/internal/rest"
	service "gitlab.com/jigsawcorp/log3900/internal/services"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/internal/services/healthcheck"
	"gitlab.com/jigsawcorp/log3900/internal/services/logger"
	"gitlab.com/jigsawcorp/log3900/internal/services/messenger"
	"gitlab.com/jigsawcorp/log3900/internal/services/router"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/model"
	"gitlab.com/jigsawcorp/log3900/pkg/graceful"
)

func main() {
	config.Init()
	model.DBConnect()
	geometry.InitTable()
	language.Init()

	restServer := &rest.Server{}
	socketServer := &socket.Server{}
	socket.RegisterBroadcast()

	graceful.Register(restServer.Shutdown, "REST server")
	graceful.Register(socketServer.Shutdown, "Socket server")
	graceful.Register(service.ShutdownAll, "Services") //Shutdown all services
	graceful.Register(model.DBClose, "Database")

	handleGraceful := graceful.ListenSIG()
	cbroadcast.NonBlockingBuffer(lockingBroadcast)

	registerServices()

	log.Printf("Server is starting jobs!")
	handleRest := make(chan bool)
	go func() {
		restServer.Initialize()
		restServer.Run(fmt.Sprintf("%s:%s", viper.GetString("rest.address"), viper.GetString("rest.port")))
		handleRest <- true
	}()

	handleSocket := make(chan bool)
	go func() {
		socketServer.StartListening(fmt.Sprintf("%s:%s", viper.GetString("socket.address"), viper.GetString("socket.port")))
		handleSocket <- true
	}()

	service.StartAll()

	<-handleRest
	<-handleSocket
	<-handleGraceful

}

func registerServices() {
	service.Add(&messenger.Messenger{})
	service.Add(&router.Router{})
	service.Add(&logger.Logger{})
	service.Add(&auth.Auth{})
	service.Add(&healthcheck.HealthCheck{})
	service.Add(&redisservice.RedisService{})
	service.Add(&potrace.Potrace{})
	service.Add(&drawing.Drawing{})
	service.Add(&lobby.Lobby{})
	service.Add(&match.Service{})
	service.Add(&virtualplayer.VirtualPlayer{})
	service.Add(&stats.Stats{})
}

func lockingBroadcast(name string) {
	log.Printf("[DeadlockWatchdog] -> Channel %s is currently blocked", name)
}
