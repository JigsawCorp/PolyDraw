package drawing

import (
	"log"

	"github.com/google/uuid"
	"github.com/tevino/abool"
	service "gitlab.com/jigsawcorp/log3900/internal/services"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"
)

//Drawing service used to route the packets and send the preview
type Drawing struct {
	preview     cbroadcast.Channel
	stopPreview cbroadcast.Channel

	strokeChunk cbroadcast.Channel
	drawStart   cbroadcast.Channel
	drawEnd     cbroadcast.Channel
	drawErase   cbroadcast.Channel

	router Router

	shutdown chan bool
}

//Init the drawing service
func (d *Drawing) Init() {
	d.shutdown = make(chan bool)
	pendingPreviews = make(map[uuid.UUID]*abool.AtomicBool)
	d.subscribe()

	d.router.Init()
}

//Start the drawing service
func (d *Drawing) Start() {
	log.Println("[Drawing] -> Starting service")
	go d.listen()
}

//Shutdown the drawing service
func (d *Drawing) Shutdown() {
	log.Println("[Drawing] -> Closing service")
	close(d.shutdown)
}

func (d *Drawing) listen() {
	defer service.Closed()

	for {
		select {
		case data := <-d.preview:
			log.Println("[Drawing] -> Receives preview request")
			if message, ok := data.(socket.RawMessageReceived); ok {
				//Start a new function to handle the connection
				go d.handlePreview(message)
			}
		case data := <-d.stopPreview:
			if message, ok := data.(socket.RawMessageReceived); ok {
				//Start a new function to handle the connection
				go stopPreview(message)
			}
		case data := <-d.strokeChunk:
			if message, ok := data.(socket.RawMessageReceived); ok {
				d.router.Route(&message)
			}
		case data := <-d.drawStart:
			if message, ok := data.(socket.RawMessageReceived); ok {
				d.router.Route(&message)
			}
		case data := <-d.drawEnd:
			if message, ok := data.(socket.RawMessageReceived); ok {
				d.router.Route(&message)
			}
		case data := <-d.drawErase:
			if message, ok := data.(socket.RawMessageReceived); ok {
				d.router.Route(&message)
			}
		case <-d.shutdown:
			return
		}
	}

}

func (d *Drawing) subscribe() {
	d.preview, _ = cbroadcast.Subscribe(BPreview)
	d.stopPreview, _ = cbroadcast.Subscribe(BStopPreview)
	d.strokeChunk, _ = cbroadcast.Subscribe(BStrokeChunk)
	d.drawStart, _ = cbroadcast.Subscribe(BDrawStart)
	d.drawEnd, _ = cbroadcast.Subscribe(BDrawEnd)
	d.drawErase, _ = cbroadcast.Subscribe(BDrawErase)
}
