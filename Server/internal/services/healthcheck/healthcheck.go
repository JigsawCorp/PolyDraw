package healthcheck

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	service "gitlab.com/jigsawcorp/log3900/internal/services"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"
)

//HealthCheck service used to debug the message received by the server
type HealthCheck struct {
	newConnection    cbroadcast.Channel
	checkReceived    cbroadcast.Channel
	connectionClosed cbroadcast.Channel

	closing     *int32
	m           sync.RWMutex
	connections map[uuid.UUID]bool
	shutdown    chan bool
}

//Init the messenger service
func (h *HealthCheck) Init() {
	h.shutdown = make(chan bool)
	h.closing = new(int32)
	*h.closing = 0
	h.connections = make(map[uuid.UUID]bool)
	h.subscribe()
}

//Start the messenger service
func (h *HealthCheck) Start() {
	log.Println("[Healthcheck] -> Starting service")
	go h.listen()
}

//Shutdown the messenger service
func (h *HealthCheck) Shutdown() {
	log.Println("[Healthcheck] -> Closing service")
	close(h.shutdown)
}

func (h *HealthCheck) listen() {
	defer service.Closed()

	for {
		select {
		case data := <-h.newConnection:
			if socketID, ok := data.(uuid.UUID); ok {
				//Start a new function to handle the connection
				go h.handleNewHost(socketID)
			}
		case data := <-h.connectionClosed:
			if socketID, ok := data.(uuid.UUID); ok {
				go h.handleCloseHost(socketID)
			}
		case data := <-h.checkReceived:
			if message, ok := data.(socket.RawMessageReceived); ok {
				h.handleReception(message)
			}
		case <-h.shutdown:
			atomic.StoreInt32(h.closing, 1)
			return
		}
	}

}

func (h *HealthCheck) subscribe() {
	h.newConnection, _ = cbroadcast.Subscribe(socket.BSocketAuthConnected)
	h.connectionClosed, _ = cbroadcast.Subscribe(socket.BSocketAuthCloseClient)
	h.checkReceived, _ = cbroadcast.Subscribe(BReceived)
}
