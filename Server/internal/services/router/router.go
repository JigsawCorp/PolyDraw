package router

import (
	"log"

	"github.com/google/uuid"
	service "gitlab.com/jigsawcorp/log3900/internal/services"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"
)

//Router the router service
type Router struct {
	receiveChan       cbroadcast.Channel
	closingSocketChan cbroadcast.Channel
	closingChan       cbroadcast.Channel
	shutdown          chan bool
	routes            map[byte]string
}

//Init the router must be called before start
func (r *Router) Init() {
	r.shutdown = make(chan bool)
	r.routes = make(map[byte]string)
	r.subscribe()

	r.routing()
}

// Register any broadcast for the router. not used for this service
func (r *Router) Register() {

}

//Start the router service
func (r *Router) Start() {
	log.Println("[SRouter] -> Starting service")
	go r.listen()
}

//Shutdown the router service
func (r *Router) Shutdown() {
	log.Println("[SRouter] -> Closing service")
	close(r.shutdown)
}

func (r *Router) listen() {
	defer service.Closed()

	for {
		select {
		case data := <-r.receiveChan:
			message, ok := data.(socket.RawMessageReceived)

			if ok {
				//Route the message to the correct service if authenticated
				if auth.IsAuthenticated(message) {
					r.route(message)
				} else {
					log.Printf("[SRouter] -> Message not authenticated!")
				}
			}
		case data := <-r.closingChan:
			socketID, ok := data.(uuid.UUID)

			if ok && auth.IsAuthenticatedQuick(socketID) {
				cbroadcast.Broadcast(socket.BSocketAuthCloseClient, socketID)
			}

		case <-r.shutdown:
			return
		}
	}

}

//route the message to the correct service
func (r *Router) route(message socket.RawMessageReceived) {
	if message.Payload.MessageType != 0 {
		if broadcast, ok := r.routes[message.Payload.MessageType]; ok {
			cbroadcast.Broadcast(broadcast, message)
		} else {
			log.Printf("[SRouter] -> No route for %d", message.Payload.MessageType)
		}
	}
}

//handle a route that the router needs to handle
func (r *Router) handle(messageType int, broadcastName string) {
	if messageType > 255 || messageType < 0 {
		panic("Can not have a message type bellow 0 or over 255")
	}
	r.routes[byte(messageType)] = broadcastName
}

func (r *Router) subscribe() {
	r.receiveChan, _ = cbroadcast.Subscribe(socket.BSocketReceive)
	r.closingSocketChan, _ = cbroadcast.Subscribe(socket.BSocketClose)
	r.closingChan, _ = cbroadcast.Subscribe(socket.BSocketCloseClient)
}
