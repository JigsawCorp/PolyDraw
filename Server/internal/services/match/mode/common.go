package mode

import (
	"context"
	"log"
	"time"

	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"

	"gitlab.com/jigsawcorp/log3900/internal/services/virtualplayer"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/tevino/abool"
	match2 "gitlab.com/jigsawcorp/log3900/internal/match"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/model"
	"golang.org/x/sync/semaphore"
)

const imageDuration = 60000

type players struct {
	socketID uuid.UUID
	userID   uuid.UUID
	Username string
	Order    int
	IsCPU    bool
}

type base struct {
	readyMatch  *semaphore.Weighted
	readyOnce   map[uuid.UUID]bool
	players     []players
	connections map[uuid.UUID]*players
	info        model.Group
	wordHistory map[string]bool
	isRunning   bool

	nbWaitingResponses int64
	waitingResponse    *semaphore.Weighted
	cancelWait         func()
	funcSyncPlayer     func()
	receivingGuesses   *abool.AtomicBool
	timeImage          int64
}

func (b *base) init(connections []uuid.UUID, info model.Group) {

	b.players = make([]players, len(connections)+info.VirtualPlayers)
	b.connections = make(map[uuid.UUID]*players, len(connections))
	b.wordHistory = make(map[string]bool)

	b.receivingGuesses = abool.New()
	b.timeImage = imageDuration

	for i := range connections {
		socketID := connections[i]
		userID, _ := auth.GetUserID(socketID)
		//Find the user data in the game info
		var user *model.User
		for j := range info.Users {
			if info.Users[j].ID == userID {
				user = info.Users[j]
			}
		}
		if user != nil && userID != uuid.Nil {
			b.players[i] = players{
				socketID: socketID,
				userID:   userID,
				Username: user.Username,
				IsCPU:    false,
			}
			b.connections[socketID] = &b.players[i]
		}
	}

	bots := virtualplayer.GetVirtualPlayersInfo(info.ID)

	if bots != nil {
		if len(bots) == info.VirtualPlayers {
			offset := len(connections)
			for i, bot := range bots {
				b.players[offset+i] = players{
					socketID: uuid.Nil,
					userID:   bot.BotID,
					Username: bot.Username,
					IsCPU:    true,
				}
			}
		}
	}

	b.info = info

	nbConnections := int64(len(b.connections))
	b.readyMatch = semaphore.NewWeighted(nbConnections)
	b.readyMatch.TryAcquire(nbConnections)

	b.readyOnce = make(map[uuid.UUID]bool)
	for i := range b.connections {
		b.readyOnce[b.connections[i].socketID] = false
	}
	log.Printf("[Match] -> Init match %s", b.info.ID)
}

//broadcast messages to all users not in parallel
func (b *base) broadcast(message *socket.RawMessage) {
	for i := range b.connections {
		socket.SendQueueMessageSocketID(*message, b.connections[i].socketID)
	}
	log.Printf("[Match] -> Message %d broadcasted, Match: %s", message.MessageType, b.info.ID)
}

//Wait for all the clients to be ready
func (b *base) waitForPlayers() bool {
	log.Printf("[Match] Waiting for all the players to register, match: %s", b.info.ID)
	nbConnections := int64(len(b.connections))

	cnt := context.Background()
	var cancel func()
	cnt, cancel = context.WithTimeout(cnt, time.Second*5)
	if b.readyMatch.Acquire(cnt, nbConnections) != nil {
		cancel()

		cancelMessage := socket.RawMessage{}
		cancelMessage.ParseMessagePack(byte(socket.MessageType.GameCancel), GameCancel{
			Type: 1,
		})
		b.broadcast(&cancelMessage)
		cbroadcast.Broadcast(match2.BGameEnds, b.info.ID)
		return false
	}
	cancel()
	//Send a message to all the clients to advise them that the game is starting
	message := socket.RawMessage{}
	message.MessageType = byte(socket.MessageType.GameStarting)
	b.broadcast(&message)
	return true
}

func (b *base) ready(socketID uuid.UUID) {
	if !b.readyOnce[socketID] {
		b.readyMatch.Release(1)
		b.readyOnce[socketID] = true
	}
}

//getPlayers used to return players needs to be exported by the implementing struct
func (b *base) getPlayers() []match2.Player {
	players := make([]match2.Player, len(b.players))
	for i := range b.players {
		players[i].ID = b.players[i].userID
		players[i].Username = b.players[i].Username
		players[i].IsCPU = b.players[i].IsCPU
	}
	return players
}

func (b *base) findGame() *model.Game {
	word := ""
	watchDog := 0

	var game model.Game
	for word == "" {
		var count int
		model.DB().Model(&game).Where("difficulty = ? and language = ?", b.info.Difficulty, b.info.Language).Count(&count)
		model.DB().Preload("Image").Preload("Hints").Joins("left join game_images on game_images.game_id = games.id").Where("difficulty = ? and language = ? and game_images.svg_file IS NOT NULL", b.info.Difficulty, b.info.Language).Order(gorm.Expr("random()")).First(&game)
		if count == 0 {
			return &game
		}

		if game.ID != uuid.Nil && len(game.Hints) > 0 && game.Image != nil {
			if _, inList := b.wordHistory[word]; !inList || watchDog >= 10 {
				//Add the word to the list so it does not come up again.
				word = game.Word
				b.wordHistory[word] = true
				return &game
			}
		}
		watchDog++
	}
	return &game
}

//waitTimeout wait for the guesses
func (b *base) waitTimeout() bool {
	c := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				//Send an update to the clients
				b.funcSyncPlayer()
			case <-c:
				return
			}
		}
	}()

	cnt := context.Background()
	cnt, b.cancelWait = context.WithTimeout(cnt, time.Millisecond*time.Duration(b.timeImage))
	err := b.waitingResponse.Acquire(cnt, b.nbWaitingResponses)
	b.cancelWait()

	close(c)

	if err == nil {
		b.receivingGuesses.UnSet()
		return false // completed normally
	}

	b.receivingGuesses.UnSet()
	return true // timed out
}

//IsRunning returns if the game is running
func (b *base) IsRunning() bool {
	return b.isRunning
}
