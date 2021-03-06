package messenger

import (
	"log"
	"strings"
	"sync"
	"time"

	"gitlab.com/jigsawcorp/log3900/internal/language"

	"gitlab.com/jigsawcorp/log3900/internal/match"
	match2 "gitlab.com/jigsawcorp/log3900/internal/match"
	"gitlab.com/jigsawcorp/log3900/pkg/cbroadcast"

	"github.com/google/uuid"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/model"
)

type handler struct {
	channelsConnections map[uuid.UUID]map[uuid.UUID]bool //channelID - socketID
	mutex               sync.RWMutex
}

func (h *handler) createGroupChannel(group *model.Group) (uuid.UUID, socket.RawMessage) {
	channel := model.ChatChannel{
		Name:       "Game",
		GroupID:    group.ID,
		IsGameChat: true,
	}
	model.DB().Create(&channel)

	cbroadcast.Broadcast(match.BChatNew, match2.ChatNew{
		MatchID: group.ID,
		ChatID:  channel.ID,
	})
	h.mutex.Lock()
	//Init the hashmap for the connections
	h.channelsConnections[channel.ID] = make(map[uuid.UUID]bool)
	h.mutex.Unlock()
	//Create request
	response := ChannelCreateResponse{
		ChannelName: channel.Name,
		ChannelID:   channel.ID.String(),
		Username:    "host",
		UserID:      uuid.Nil.String(),
		Timestamp:   time.Now().Unix(),
		IsGame:      true,
	}
	rawMessage := socket.RawMessage{}
	if rawMessage.ParseMessagePack(byte(socket.MessageType.UserCreateChannel), response) != nil {
		log.Printf("[Messenger] -> Create: Can't pack message. Dropping packet!")
		return channel.ID, socket.RawMessage{}
	}
	log.Printf("[Messenger] -> Create: channel %s created", channel.Name)

	return channel.ID, rawMessage
}

func (h *handler) deleteGroupChannel(group *model.Group) {
	var channel model.ChatChannel
	model.DB().Where("group_id = ?", group.ID).First(&channel)
	if channel.ID != uuid.Nil {
		//Create a destroy message
		destroyMessage := ChannelDestroyResponse{
			UserID:    uuid.Nil.String(),
			Username:  "host",
			ChannelID: channel.ID.String(),
			Timestamp: time.Now().Unix(),
		}
		rawMessage := socket.RawMessage{}
		if rawMessage.ParseMessagePack(byte(socket.MessageType.UserDestroyedChannel), destroyMessage) != nil {
			log.Printf("[Messenger] -> Destroy: Can't pack message. Dropping packet!")
			return
		}
		h.mutex.Lock()
		h.broadcast(channel.ID, rawMessage)

		delete(h.channelsConnections, channel.ID)
		h.mutex.Unlock()
		model.DB().Delete(&channel)
	}

}

func (h *handler) quitChannel(socketID uuid.UUID, channelID uuid.UUID) {
	//Check if channel exists
	channel := model.ChatChannel{}
	model.DB().Model(&channel).Related(&model.User{}, "Users")
	model.DB().Preload("Users").Where("id = ?", channelID).First(&channel)
	if channel.ID != uuid.Nil {
		user, _ := auth.GetUser(socketID)
		h.mutex.Lock()
		if _, ok := h.channelsConnections[channel.ID][socketID]; ok {
			model.DB().Model(&channel).Association("Users").Delete(user)

			//Create a quit message
			quitResponse := ChannelLeaveResponse{
				UserID:    user.ID.String(),
				Username:  user.Username,
				ChannelID: channel.ID.String(),
				Timestamp: time.Now().Unix(),
			}
			rawMessage := socket.RawMessage{}
			if rawMessage.ParseMessagePack(byte(socket.MessageType.UserLeftChannel), quitResponse) != nil {
				h.mutex.Unlock()
				log.Printf("[Messenger] -> Quit: Can't pack message. Dropping packet!")
				return
			}

			h.broadcast(channel.ID, rawMessage)
			delete(h.channelsConnections[channelID], socketID)
			h.mutex.Unlock()
			log.Printf("[Messenger] -> Quit: User %s quit %s", user.ID.String(), channelID)
		} else {
			h.mutex.Unlock()
			log.Printf("[Messenger] -> Quit: User is not in the channel")
		}
	} else {
		log.Printf("[Messenger] -> Quit: Invalid channel UUID, not found")
		socket.SendErrorToSocketID(socket.MessageType.LeaveChannel, 404, language.MustGetSocket("error.channelNotFound", socketID), socketID)
	}
}

func (h *handler) joinChannel(socketID uuid.UUID, channelID uuid.UUID) {
	channel := model.ChatChannel{}
	model.DB().Model(&channel).Related(&model.User{}, "Users")
	model.DB().Preload("Users").Where("id = ?", channelID).First(&channel)

	if channel.ID != uuid.Nil {
		user, _ := auth.GetUser(socketID)
		h.mutex.Lock()
		if _, ok := h.channelsConnections[channel.ID][socketID]; !ok {
			joinServer := ChannelJoin{
				UserID:    user.ID.String(),
				Username:  user.Username,
				ChannelID: channel.ID.String(),
				Timestamp: time.Now().Unix(),
			}

			rawMessage := socket.RawMessage{}
			if rawMessage.ParseMessagePack(byte(socket.MessageType.UserJoinedChannel), joinServer) != nil {
				h.mutex.Unlock()
				log.Printf("[Messenger] -> Join: Can't pack message. Dropping packet!")
				return
			}

			//We can join the channel
			model.DB().Model(&channel).Association("Users").Append(user)
			h.channelsConnections[channel.ID][socketID] = true

			h.broadcast(channel.ID, rawMessage)
			h.mutex.Unlock()
			log.Printf("[Messenger] -> Join: User %s join %s", user.ID.String(), channelID)
		} else {
			h.mutex.Unlock()
			log.Printf("[Messenger] -> Join: User is already joined to the channel")
			socket.SendErrorToSocketID(socket.MessageType.JoinChannel, 409, language.MustGetSocket("error.userJoinChannel", socketID), socketID)
		}
	} else {
		log.Printf("[Messenger] -> Join: Channel UUID not found, %s", channelID.String())
		socket.SendErrorToSocketID(socket.MessageType.JoinChannel, 404, language.MustGetSocket("error.channelNotFound", socketID), socketID)
	}
}

func (h *handler) broadcast(chanelID uuid.UUID, message socket.RawMessage) {
	for socketID := range h.channelsConnections[chanelID] {
		socket.SendQueueMessageSocketID(message, socketID)
	}
}

func (h *handler) init() {
	h.channelsConnections = map[uuid.UUID]map[uuid.UUID]bool{}
	h.channelsConnections[uuid.Nil] = make(map[uuid.UUID]bool)
	setInstance(h)

	var channels []model.ChatChannel
	model.DB().Find(&channels)

	for i := range channels {
		h.channelsConnections[channels[i].ID] = make(map[uuid.UUID]bool)
	}
}
func (h *handler) handleMessgeSent(message socket.RawMessageReceived) {
	var messageParsed MessageSent
	timestamp := time.Now().Unix()
	if message.Payload.DecodeMessagePack(&messageParsed) == nil {
		//Send to all other connected users
		user, err := auth.GetUser(message.SocketID)
		if err != nil {
			log.Printf("[Messenger] -> %s", err)
		}
		channelID, err := uuid.Parse(messageParsed.ChannelID)
		if err == nil {
			h.mutex.RLock()
			if _, ok := h.channelsConnections[channelID][message.SocketID]; ok {
				messageToFoward := MessageReceived{
					ChannelID: messageParsed.ChannelID,
					UserID:    user.ID.String(),
					Username:  user.Username,
					Message:   messageParsed.Message,
					Timestamp: timestamp,
				}
				rawMessage := socket.RawMessage{}
				if rawMessage.ParseMessagePack(byte(socket.MessageType.MessageReceived), messageToFoward) != nil {
					h.mutex.RUnlock()
					log.Printf("[Messenger] -> Receive: Can't pack message. Dropping packet!")
					return
				}
				h.broadcast(channelID, rawMessage)
				h.mutex.RUnlock()
				log.Printf("[Messenger] -> Receive: \"%s\" Username: \"%s\" ChannelID: %s", messageParsed.Message, user.Username, messageParsed.ChannelID)
				model.AddMessage(messageParsed.Message, channelID, user.ID, timestamp)
			} else {
				h.mutex.RUnlock()
				log.Printf("[Messenger] -> Receive: The user needs to join the channel first. Dropping packet!")
				socket.SendErrorToSocketID(socket.MessageType.MessageSent, 409, language.MustGetSocket("error.userJoinFirst", message.SocketID), message.SocketID)
			}
		} else {
			log.Printf("[Messenger] -> Receive: Invalid channel ID. Dropping packet!")
			socket.SendErrorToSocketID(socket.MessageType.MessageSent, 404, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
		}
	} else {
		log.Printf("[Messenger] -> Receive: Wrong data format. Dropping packet!")
		socket.SendErrorToSocketID(socket.MessageType.MessageSent, 404, "Wrong data format.", message.SocketID)
	}
}

func (h *handler) handleCreateChannel(message socket.RawMessageReceived) {
	var channelParsed ChannelCreateRequest
	timestamp := time.Now().Unix()
	if message.Payload.DecodeMessagePack(&channelParsed) == nil {
		name := channelParsed.ChannelName
		if strings.TrimSpace(name) != "" &&
			name != "Général" &&
			name != "General" &&
			name != "Game" &&
			name != "Active Match" &&
			name != "Partie en cours" {
			user, err := auth.GetUser(message.SocketID)
			if err == nil {
				//Check if channel already exists
				var count int64
				model.DB().Model(&model.ChatChannel{}).Where("name = ?", name).Count(&count)
				if count == 0 {
					channel := model.ChatChannel{
						Name: name,
					}
					model.DB().Create(&channel)

					h.mutex.Lock()
					//Init the hashmap for the connections
					h.channelsConnections[channel.ID] = make(map[uuid.UUID]bool)
					//Create request
					response := ChannelCreateResponse{
						ChannelName: name,
						ChannelID:   channel.ID.String(),
						Username:    user.Username,
						UserID:      user.ID.String(),
						IsGame:      false,
						Timestamp:   timestamp,
					}
					rawMessage := socket.RawMessage{}
					if rawMessage.ParseMessagePack(byte(socket.MessageType.UserCreateChannel), response) != nil {
						h.mutex.Unlock()
						log.Printf("[Messenger] -> Create: Can't pack message. Dropping packet!")
						return
					}

					h.broadcast(uuid.Nil, rawMessage)
					h.mutex.Unlock()
					log.Printf("[Messenger] -> Create: channel %s created", channelParsed.ChannelName)
				} else {
					log.Printf("[Messenger] -> Create: Channel already exists. Dropping packet!")
					socket.SendErrorToSocketID(socket.MessageType.CreateChannel, 409, language.MustGetSocket("error.channelExists", message.SocketID), message.SocketID)
				}
			} else {
				log.Printf("[Messenger] -> Create: Can't find user. Dropping packet!")
				socket.SendErrorToSocketID(socket.MessageType.CreateChannel, 404, language.MustGetSocket("error.userNotFound", message.SocketID), message.SocketID)
			}
		} else {
			log.Printf("[Messenger] -> Create: Invalid channel name. Dropping packet!")
			socket.SendErrorToSocketID(socket.MessageType.CreateChannel, 400, language.MustGetSocket("error.channelInvalidName", message.SocketID), message.SocketID)
		}
	} else {
		log.Printf("[Messenger] -> Create: Invalid channel decoding. Dropping packet!")
		socket.SendErrorToSocketID(socket.MessageType.CreateChannel, 400, "Invalid channel decoding.", message.SocketID)
	}
}

func (h *handler) handleJoinChannel(message socket.RawMessageReceived) {
	if message.Payload.Length == 16 {
		channelID, err := uuid.FromBytes(message.Payload.Bytes)
		if err == nil {
			h.joinChannel(message.SocketID, channelID)
		} else {
			log.Printf("[Messenger] -> Join: Invalid channel UUID")
			socket.SendErrorToSocketID(socket.MessageType.JoinChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
		}
	} else {
		log.Printf("[Messenger] -> Join: Invalid channel UUID")
		socket.SendErrorToSocketID(socket.MessageType.JoinChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
	}
}

func (h *handler) handleQuitChannel(message socket.RawMessageReceived) {
	if message.Payload.Length == 16 {
		channelID, err := uuid.FromBytes(message.Payload.Bytes)
		if err == nil {
			h.quitChannel(message.SocketID, channelID)
		} else {
			log.Printf("[Messenger] -> Quit: Invalid channel UUID")
			socket.SendErrorToSocketID(socket.MessageType.LeaveChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
		}
	} else {
		log.Printf("[Messenger] -> Quit: Invalid channel UUID")
		socket.SendErrorToSocketID(socket.MessageType.LeaveChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
	}
}

func (h *handler) handleDestroyChannel(message socket.RawMessageReceived) {
	if message.Payload.Length == 16 {
		channelID, err := uuid.FromBytes(message.Payload.Bytes)
		if err == nil {
			//Check if channel exists
			channel := model.ChatChannel{}
			model.DB().Model(&channel).Related(&model.User{}, "Users")
			model.DB().Preload("Users").Where("id = ?", channelID).First(&channel)

			if channel.ID != uuid.Nil {
				user, _ := auth.GetUser(message.SocketID)

				h.mutex.Lock()
				delete(h.channelsConnections, channel.ID)
				h.mutex.Unlock()

				model.DB().Model(&channel).Delete(&channel)

				//Create a destroy message
				destroyMessage := ChannelDestroyResponse{
					UserID:    user.ID.String(),
					Username:  user.Username,
					ChannelID: channel.ID.String(),
					Timestamp: time.Now().Unix(),
				}
				rawMessage := socket.RawMessage{}
				if rawMessage.ParseMessagePack(byte(socket.MessageType.UserDestroyedChannel), destroyMessage) != nil {
					log.Printf("[Messenger] -> Destroy: Can't pack message. Dropping packet!")
					return
				}
				h.mutex.RLock()
				h.broadcast(uuid.Nil, rawMessage)
				h.mutex.RUnlock()

				log.Printf("[Messenger] -> Destroy: Removed channel %s", channelID)
			} else {
				log.Printf("[Messenger] -> Destroy: Invalid channel UUID, not found")
				socket.SendErrorToSocketID(socket.MessageType.DestroyChannel, 404, language.MustGetSocket("error.channelNotFound", message.SocketID), message.SocketID)
			}
		} else {
			log.Printf("[Messenger] -> Destroy: Invalid channel UUID")
			socket.SendErrorToSocketID(socket.MessageType.DestroyChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
		}
	} else {
		log.Printf("[Messenger] -> Destroy: Invalid channel UUID")
		socket.SendErrorToSocketID(socket.MessageType.DestroyChannel, 400, language.MustGetSocket("error.channelInvalidUUID", message.SocketID), message.SocketID)
	}
}

func (h *handler) handleConnect(socketID uuid.UUID) {
	h.mutex.Lock()
	h.channelsConnections[uuid.Nil][socketID] = true
	h.mutex.Unlock()

	user, _ := auth.GetUser(socketID)

	var channels []model.ChatChannel

	model.DB().Joins("left join chat_channel_membership on chat_channel_membership.chat_channel_id = chat_channels.id").Where("chat_channel_membership.user_id = ?", user.ID).Find(&channels)
	joinServer := ChannelJoin{
		UserID:    user.ID.String(),
		Username:  user.Username,
		ChannelID: uuid.Nil.String(),
		Timestamp: time.Now().Unix(),
	}

	rawMessage := socket.RawMessage{}
	if rawMessage.ParseMessagePack(byte(socket.MessageType.UserJoinedChannel), joinServer) != nil {
		log.Printf("[Messenger] -> Connect: Can't pack message. Dropping packet!")
		return
	}

	h.mutex.Lock()
	h.broadcast(uuid.Nil, rawMessage)

	//Update the cache
	for _, channel := range channels {
		h.channelsConnections[channel.ID][socketID] = true
	}
	h.mutex.Unlock()
}

func (h *handler) handleDisconnect(socketID uuid.UUID) {
	h.mutex.Lock()
	delete(h.channelsConnections[uuid.Nil], socketID)
	h.mutex.Unlock()

	user, _ := auth.GetUser(socketID)

	var channels []model.ChatChannel
	model.DB().Joins("left join chat_channel_membership on chat_channel_membership.chat_channel_id = chat_channels.id").Where("chat_channel_membership.user_id = ?", user.ID).Find(&channels)

	leaveServer := ChannelJoin{
		UserID:    user.ID.String(),
		Username:  user.Username,
		ChannelID: uuid.Nil.String(),
		Timestamp: time.Now().Unix(),
	}

	rawMessage := socket.RawMessage{}
	if rawMessage.ParseMessagePack(byte(socket.MessageType.UserLeftChannel), leaveServer) != nil {
		log.Printf("[Messenger] -> Disconnect: Can't pack message. Dropping packet!")
		return
	}
	h.mutex.Lock()
	h.broadcast(uuid.Nil, rawMessage)

	//Update the cache
	for _, channel := range channels {
		if !channel.IsGameChat {
			delete(h.channelsConnections[channel.ID], socketID)
		}
	}
	h.mutex.Unlock()
}

func (h *handler) handleBotMessage(message MessageReceived) {
	channelID, err := uuid.Parse(message.ChannelID)
	if err == nil {
		rawMessage := socket.RawMessage{}
		if rawMessage.ParseMessagePack(byte(socket.MessageType.MessageReceived), message) != nil {
			log.Printf("[Messenger] -> Receive: Can't pack message. Dropping packet!")
			return
		}
		h.mutex.RLock()
		for k := range h.channelsConnections[channelID] {
			// Send message to the socket in async way
			socket.SendQueueMessageSocketID(rawMessage, k)
		}
		h.mutex.RUnlock()
		log.Printf("[Messenger] -> Receive: \"%s\" Username: \"%s\" ChannelID: %s", message.Message, message.Username, message.ChannelID)
		botID, err2 := uuid.Parse(message.UserID)
		if err2 == nil {
			model.AddMessage(message.Message, channelID, botID, time.Now().Unix())
		} else {
			log.Printf("[Messenger] -> Receive: Invalid bot ID. Dropping packet!")
		}
	} else {
		log.Printf("[Messenger] -> Receive: Invalid channel ID. Dropping packet!")
	}

}
