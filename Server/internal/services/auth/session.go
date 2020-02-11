package auth

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/model"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"
)

//Session represents a session with authentification
type Session struct {
	UserID uuid.UUID
	Token  string
}

var mutex sync.Mutex
var tokenAvailable map[string]uuid.UUID
var sessionCache map[uuid.UUID]bool

//TODO move to the service and init once
func initTokenAvailable() {
	if tokenAvailable == nil {
		tokenAvailable = make(map[string]uuid.UUID)
	}
	if sessionCache == nil {
		sessionCache = make(map[uuid.UUID]bool)
	}
}

//Register the token that the client can use for authentification
func Register(token string, userID uuid.UUID) {
	defer mutex.Unlock()
	mutex.Lock()

	initTokenAvailable()
	tokenAvailable[token] = userID
	log.Printf("[Auth] -> Registering user ID: %s", userID)

	//TODO cleanup of unused tokens
}

//UnRegisterSocket removes the session from the socketID
func UnRegisterSocket(socketID uuid.UUID) {
	defer mutex.Unlock()
	mutex.Lock()

	var session model.Session
	if removingSessions.IsSet() {
		return
	}
	model.DB().Where("socket_id = ?", socketID).First(&session)

	if session.ID != uuid.Nil {
		token := session.SessionToken

		delete(tokenAvailable, token)
		delete(sessionCache, socketID)

		model.DB().Delete(&session) //Remove the session
	}
}

//UnRegisterUser removes the session from the userID
func UnRegisterUser(userID uuid.UUID) {
	defer mutex.Unlock()
	mutex.Lock()

	var session model.Session
	if removingSessions.IsSet() {
		return
	}
	model.DB().Where("user_id = ?", userID).First(&session)

	if session.ID != uuid.Nil {
		token := session.SessionToken

		delete(tokenAvailable, token)
		delete(sessionCache, session.SocketID)

		model.DB().Delete(&session) //Remove the session
	}
}

//IsTokenAvailable makes sure that the token does not already exist in the session table
func IsTokenAvailable(token string) bool {
	defer mutex.Unlock()
	mutex.Lock()

	initTokenAvailable()
	_, ok := tokenAvailable[token]
	return !ok
}

//IsAuthenticated returns if the user is authenticated by the rest API. Parse the message if it is the correct type
func IsAuthenticated(messageReceived socket.RawMessageReceived) bool {
	defer mutex.Unlock()
	mutex.Lock()

	initTokenAvailable()
	if messageReceived.Payload.MessageType == byte(socket.MessageType.ServerConnection) {

		bytes := messageReceived.Payload.Bytes
		token := string(bytes)

		if userID, ok := tokenAvailable[token]; ok {

			if HasUserSession(userID) {
				log.Printf("[Auth] -> Connection already exists dropping %s", messageReceived.SocketID)
				sendAuthResponse(false, messageReceived.SocketID)
				return false
			}

			model.DB().Create(&model.Session{
				UserID:       userID,
				SessionToken: token,
				SocketID:     messageReceived.SocketID,
			})

			sessionCache[messageReceived.SocketID] = true //Set the value in the cache so pacquets are routed fast
			log.Printf("[Auth] -> Connection made %s", messageReceived.SocketID)
			sendAuthResponse(true, messageReceived.SocketID)

			cbroadcast.Broadcast(socket.BSocketAuthConnected, messageReceived.SocketID) //Broadcast only when the auth is connected
			return true
		}
		sendAuthResponse(false, messageReceived.SocketID)
		return false
	}

	_, ok := sessionCache[messageReceived.SocketID]
	if !ok {
		//Send a fail message any time if the socket is not correct or if the user was forced disconnected
		sendAuthResponse(false, messageReceived.SocketID)
	}
	return ok
}

//GetUserID returns the user associated with a session
func GetUserID(socketID uuid.UUID) (model.User, error) {
	var session model.Session
	var user model.User
	model.DB().Where("socket_id = ?", socketID).First(&session)
	model.DB().Model(&session).Related(&user)
	if user.ID != uuid.Nil {
		return user, nil
	}

	return model.User{}, fmt.Errorf("No user is associated with this connection")
}

//HasUserSession returns true if the user has a session
func HasUserSession(userID uuid.UUID) bool {
	var session model.Session
	model.DB().Where("user_id = ?", userID).First(&session)

	return session.ID != uuid.Nil
}

//HasUserToken returns true if the user has already a token issued
func HasUserToken(userID uuid.UUID) (bool, string) {
	defer mutex.Unlock()
	mutex.Lock()

	initTokenAvailable()
	for k, v := range tokenAvailable {
		if v == userID {
			return true, k
		}
	}
	return false, ""
}

//IsAuthenticatedQuick returns only if the socket is in the cache of authenticated sockets
func IsAuthenticatedQuick(socketID uuid.UUID) bool {
	defer mutex.Unlock()
	mutex.Lock()

	initTokenAvailable()

	_, ok := sessionCache[socketID]
	return ok
}

func sendAuthResponse(response bool, socketID uuid.UUID) {
	message := socket.RawMessage{}
	message.MessageType = byte(socket.MessageType.ServerConnectionResponse)
	message.Length = 1
	if response {
		message.Bytes = []byte{0x01}
	} else {
		message.Bytes = []byte{0x00}
	}
	socket.SendRawMessageToSocketID(message, socketID)
}