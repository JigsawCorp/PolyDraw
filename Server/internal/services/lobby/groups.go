package lobby

import (
	"fmt"
	"log"
	"sync"

	"gitlab.com/jigsawcorp/log3900/internal/language"

	"gitlab.com/jigsawcorp/log3900/internal/services/messenger"
	"gitlab.com/jigsawcorp/log3900/internal/services/virtualplayer"

	"github.com/google/uuid"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/internal/services/match"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
	"gitlab.com/jigsawcorp/log3900/model"
)

type responseGen struct {
	Response bool
	Error    string
}

type responseGroup struct {
	UserID   string
	Username string
	GroupID  string
	IsCPU    bool
	IsKicked bool
}

type responsePlayer struct {
	IsCPU    bool
	ID       string
	Username string
}

type responseNewGroup struct {
	ID         string
	GroupName  string
	OwnerName  string
	OwnerID    string
	PlayersMax int
	Mode       int
	Players    []responsePlayer
	NbRound    int
	Language   int
	Difficulty int
}

type groups struct {
	mutex      sync.Mutex
	queue      map[uuid.UUID]bool        // socketID -> _
	assignment map[uuid.UUID]uuid.UUID   // socketID -> groupID
	groups     map[uuid.UUID][]uuid.UUID // groupID -> socketID
}

func (g *groups) Init() {
	g.assignment = make(map[uuid.UUID]uuid.UUID)
	g.groups = make(map[uuid.UUID][]uuid.UUID)
	g.queue = make(map[uuid.UUID]bool)
}

//RegisterSession used to register a session in the queue to be added to a groups
func (g *groups) RegisterSession(socketID uuid.UUID) {
	defer g.mutex.Unlock()

	g.mutex.Lock()
	g.queue[socketID] = true
}

//UnRegisterSession used to remove the user from the groups or the waiting list
func (g *groups) UnRegisterSession(socketID uuid.UUID) {
	defer g.mutex.Unlock()
	g.mutex.Lock()

	delete(g.queue, socketID)
	if groupID, ok := g.assignment[socketID]; ok {
		delete(g.assignment, socketID)
		g.removeSocketGroup(socketID, groupID)

		userID, err := auth.GetUserID(socketID)

		var groupDB model.Group
		model.DB().Where("id = ?", groupID).First(&groupDB)
		if groupDB.ID != uuid.Nil && err == nil {
			model.DB().Model(&groupDB).Association("Users").Delete(&model.User{Base: model.Base{ID: userID}})
			messenger.HandleQuitGroup(&groupDB, socketID)
			//If user is owner we delete the group
			if groupDB.OwnerID == userID {
				g.safeDeleteGroup(&groupDB)
			}
		}
	}
}

func (g *groups) KickUser(socketID uuid.UUID, userID uuid.UUID) {
	g.mutex.Lock()
	//Find the group and find the user socket id
	if groupID, ok := g.assignment[socketID]; ok {
		var groupDB model.Group
		model.DB().Where("id = ?", groupID).First(&groupDB)
		if groupDB.ID != uuid.Nil {
			//Make sure that we are the owner
			currentUserID, _ := auth.GetUserID(socketID)
			if groupDB.OwnerID == currentUserID {
				socketKickUser, err := auth.GetSocketID(userID)
				if err == nil {
					g.mutex.Unlock()
					g.QuitGroup(socketKickUser, true)
				} else {
					g.mutex.Unlock()
					if !g.kickVirtualPlayer(userID) {
						socket.SendErrorToSocketID(socket.MessageType.RequestKickUser, 404, language.MustGetSocket("error.userNotFound", socketID), socketID)
					}
				}
			} else {
				g.mutex.Unlock()
				socket.SendErrorToSocketID(socket.MessageType.RequestKickUser, 400, language.MustGetSocket("error.groupOwner", socketID), socketID)
			}
		} else {
			g.mutex.Unlock()
			socket.SendErrorToSocketID(socket.MessageType.RequestKickUser, 404, language.MustGetSocket("error.groupNotFound", socketID), socketID)
		}
	} else {
		g.mutex.Unlock()
		socket.SendErrorToSocketID(socket.MessageType.RequestKickUser, 400, language.MustGetSocket("error.groupMembership", socketID), socketID)
	}
}

//AddGroup add the group and send the message to all the users in queue
func (g *groups) AddGroup(group *model.Group) {
	defer g.mutex.Unlock()

	message := socket.RawMessage{}
	players := []responsePlayer{
		{
			IsCPU:    false,
			ID:       group.OwnerID.String(),
			Username: group.Owner.Username,
		},
	}
	message.ParseMessagePack(byte(socket.MessageType.ResponseGroupCreated), responseNewGroup{
		ID:         group.ID.String(),
		GroupName:  group.Name,
		OwnerName:  group.Owner.Username,
		OwnerID:    group.OwnerID.String(),
		PlayersMax: group.PlayersMax,
		Mode:       group.GameType,
		Players:    players,
		Language:   group.Language,
		NbRound:    group.NbRound,
		Difficulty: group.Difficulty,
	})
	messenger.RegisterGroup(group)
	virtualplayer.AddGroup(group.ID)
	g.mutex.Lock()
	g.groups[group.ID] = make([]uuid.UUID, 0, 4)
	for k := range g.queue {
		socket.SendQueueMessageSocketID(message, k)
	}
}

//JoinGroup used to add a user to the groups can be called in rest that's why we can avoid the db operation
func (g *groups) JoinGroup(socketID uuid.UUID, groupID uuid.UUID) {
	g.mutex.Lock()
	if _, ok := g.queue[socketID]; ok {
		//Check if groups exists
		groupDB := model.Group{}
		model.DB().Where("id = ? and status = 0", groupID).First(&groupDB)
		if groupDB.ID != uuid.Nil {

			//Is the group full ?
			if groupDB.PlayersMax-groupDB.VirtualPlayers-len(g.groups[groupID]) > 0 {
				delete(g.queue, socketID)
				g.assignment[socketID] = groupID

				g.groups[groupID] = append(g.groups[groupID], socketID)
				g.mutex.Unlock()

				//send response to client
				message := socket.RawMessage{}
				message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
					Response: true,
					Error:    "",
				})

				if socket.SendRawMessageToSocketID(message, socketID) == nil {
					//We only commit the data to the db if the message was sent successfully
					//else we will handle the error in the disconnect message
					user, _ := auth.GetUser(socketID)
					model.DB().Model(&groupDB).Association("Users").Append(&user)

					//Send a message to all the member of the group to advertise that a new user is in the group
					newUser := socket.RawMessage{}
					newUser.ParseMessagePack(byte(socket.MessageType.UserJoinedGroup), responseGroup{
						UserID:   user.ID.String(),
						Username: user.Username,
						GroupID:  groupID.String(),
						IsCPU:    false,
					})
					g.mutex.Lock()
					for i := range g.groups[groupID] {
						socket.SendQueueMessageSocketID(newUser, g.groups[groupID][i])
					}
					for k := range g.queue {
						socket.SendQueueMessageSocketID(newUser, k)
					}
					g.mutex.Unlock()

					messenger.HandleJoinGroup(&groupDB, socketID)

				}
				return
			}
			g.mutex.Unlock()
			message := socket.RawMessage{}
			message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
				Response: false,
				Error:    language.MustGetSocket("error.groupIsFull", socketID),
			})
			socket.SendQueueMessageSocketID(message, socketID)
			return

		}
		g.mutex.Unlock()

		message := socket.RawMessage{}
		message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
			Response: false,
			Error:    language.MustGetSocket("error.groupNotFound", socketID),
		})
		socket.SendQueueMessageSocketID(message, socketID)
	} else {
		g.mutex.Unlock()

		message := socket.RawMessage{}
		message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
			Response: false,
			Error:    language.MustGetSocket("error.userSingleGroup", socketID),
		})
		socket.SendQueueMessageSocketID(message, socketID)
	}
}

//QuitGroup quits the groups the user is currently in.
func (g *groups) QuitGroup(socketID uuid.UUID, forced bool) {
	g.mutex.Lock()
	if _, ok := g.assignment[socketID]; ok {
		groupID := g.assignment[socketID]

		delete(g.assignment, socketID)
		g.removeSocketGroup(socketID, groupID)
		g.queue[socketID] = true

		//Send a message to the groups member to advertise that the user quit the groups
		user, _ := auth.GetUser(socketID)
		message := socket.RawMessage{}
		message.ParseMessagePack(byte(socket.MessageType.ResponseLeaveGroup), responseGroup{
			UserID:   user.ID.String(),
			Username: user.Username,
			GroupID:  groupID.String(),
			IsKicked: forced,
			IsCPU:    false,
		})
		for i := range g.groups[groupID] {
			socket.SendQueueMessageSocketID(message, g.groups[groupID][i])
		}
		for k := range g.queue {
			socket.SendQueueMessageSocketID(message, k)
		}
		g.mutex.Unlock()

		var groupDB model.Group
		model.DB().Where("id = ?", groupID).First(&groupDB)
		model.DB().Model(&groupDB).Association("Users").Delete(&user)

		messenger.HandleQuitGroup(&groupDB, socketID)

		g.mutex.Lock()
		if user.ID == groupDB.OwnerID {
			//The owner has left the group
			g.safeDeleteGroup(&groupDB)
		}
		g.mutex.Unlock()

	} else {
		g.mutex.Unlock()
	}
}

//Set all groups in the database who are set to status waiting to abandoned
func (g *groups) CleanAllGroups() {
	var groups []model.Group
	model.DB().Model(&groups).Where("status = 0").Find(&groups)
	for i := range groups {
		groups[i].Status = 3
		model.DB().Save(groups[i])
	}
}

func (g *groups) removeSocketGroup(socketID uuid.UUID, groupID uuid.UUID) {
	for i := range g.groups[groupID] {
		if g.groups[groupID][i] == socketID {
			last := len(g.groups[groupID]) - 1

			g.groups[groupID][i] = g.groups[groupID][last]
			g.groups[groupID][last] = uuid.Nil
			g.groups[groupID] = g.groups[groupID][:last]
			return
		}
	}
}

func (g *groups) safeDeleteGroup(groupDB *model.Group) {
	//We remove the group
	uuidBytes, _ := groupDB.ID.MarshalBinary()
	message := socket.RawMessage{
		MessageType: byte(socket.MessageType.ResponseGroupRemoved),
		Length:      uint16(len(uuidBytes)),
		Bytes:       uuidBytes,
	}
	//Broadcast a message to all the users in queue
	for k := range g.queue {
		socket.SendQueueMessageSocketID(message, k)
	}

	//Broadcast a message that the group was deleted and remove them from the group
	for _, v := range g.groups[groupDB.ID] {
		socket.SendQueueMessageSocketID(message, v)
	}

	messenger.UnRegisterGroup(groupDB, g.groups[groupDB.ID])
	//Remove all the data associated with the groups
	for _, v := range g.groups[groupDB.ID] {
		delete(g.assignment, v)
		g.queue[v] = true
	}
	delete(g.groups, groupDB.ID)

	virtualplayer.RemoveGroup(groupDB.ID)
	groupDB.Status = 3
	model.DB().Save(&groupDB)
}

//StartMatch method used to create the match
func (g *groups) StartMatch(socketID uuid.UUID) {

	g.mutex.Lock()
	groupID, hasGroup := g.assignment[socketID]
	g.mutex.Unlock()

	if hasGroup {
		var groupDB model.Group
		model.DB().Model(&groupDB).Related(&model.User{}, "Users")
		model.DB().Preload("Users").Where("id = ?", groupID).First(&groupDB)

		userID, _ := auth.GetUserID(socketID)
		//Start only if the owner
		if groupDB.OwnerID == userID {
			//Check if there are enough people
			g.mutex.Lock()
			count := len(g.groups[groupID])
			g.mutex.Unlock()

			if (count >= 2 && groupDB.GameType == 0) ||
				(count == 1 && groupDB.GameType == 1) ||
				(count >= 2 && groupDB.GameType == 2) {
				if groupDB.GameType >= 1 && groupDB.VirtualPlayers < 1 {
					g.addBotGroupID(&groupDB)
				}

				//Check if there are enough games
				if groupDB.VirtualPlayers >= 1 {
					var count int
					model.DB().Model(&model.Game{}).Joins("left join game_images on game_images.game_id = games.id").Where("difficulty = ? and language = ? and game_images.svg_file IS NOT NULL", groupDB.Difficulty, groupDB.Language).Count(&count)
					if count < 10 {
						rawMessage := socket.RawMessage{}
						rawMessage.ParseMessagePack(byte(socket.MessageType.ResponseGameStart), responseGen{
							Response: false,
							Error:    fmt.Sprintf(language.MustGetSocket("error.gameMinimum", socketID), count),
						})
						socket.SendQueueMessageSocketID(rawMessage, socketID)
						return
					}
				}

				//We send the response and we pass it along to the match service
				rawMessage := socket.RawMessage{}
				rawMessage.ParseMessagePack(byte(socket.MessageType.ResponseGameStart), responseGen{
					Response: true,
				})
				g.mutex.Lock()
				for _, v := range g.groups[groupDB.ID] {
					socket.SendQueueMessageSocketID(rawMessage, v)
				}
				uuidBytes, _ := groupDB.ID.MarshalBinary()
				message := socket.RawMessage{
					MessageType: byte(socket.MessageType.ResponseGroupRemoved),
					Length:      uint16(len(uuidBytes)),
					Bytes:       uuidBytes,
				}

				//Broadcast a message to all the users in queue
				for k := range g.queue {
					socket.SendQueueMessageSocketID(message, k)
				}
				connections := g.groups[groupID][:]
				g.mutex.Unlock()

				match.UpgradeGroup(groupID, connections, groupDB)
				groupDB.Status = 1
				model.DB().Save(&groupDB)

				//change status and put all the users in the queue once they quit the game
				//Remove all the data associated with the groups
				g.mutex.Lock()
				for _, v := range g.groups[groupDB.ID] {
					delete(g.assignment, v)
					g.queue[v] = true
				}
				delete(g.groups, groupDB.ID)
				g.mutex.Unlock()
			} else {
				rawMessage := socket.RawMessage{}
				rawMessage.ParseMessagePack(byte(socket.MessageType.ResponseGameStart), responseGen{
					Response: false,
					Error:    language.MustGetSocket("error.notEnough", socketID),
				})
				socket.SendQueueMessageSocketID(rawMessage, socketID)
			}

		} else {
			rawMessage := socket.RawMessage{}
			rawMessage.ParseMessagePack(byte(socket.MessageType.ResponseGameStart), responseGen{
				Response: false,
				Error:    "Only the owner can request the game to start.",
			})
			socket.SendQueueMessageSocketID(rawMessage, socketID)
		}
	} else {
		rawMessage := socket.RawMessage{}
		rawMessage.ParseMessagePack(byte(socket.MessageType.ResponseGameStart), responseGen{
			Response: false,
			Error:    "The user is not associated with any group.",
		})
		socket.SendQueueMessageSocketID(rawMessage, socketID)
	}
}

//addBotGroupID, remove
func (g *groups) addBotGroupID(groupDB *model.Group) {
	user := model.User{IsCPU: true}
	model.AddUser(&user)

	user.Username = virtualplayer.AddVirtualPlayer(groupDB.ID, user.ID, groupDB.Language)
	log.Printf("[Lobby] -> adding bot in DB: %v", user)

	if user.Username != "" {
		//appel au joueur virtuel
		model.DB().Model(groupDB).Association("Users").Append(&user)
		groupDB.VirtualPlayers++
		model.DB().Save(groupDB)

		//Send a message to all the member of the group to advertise that a new user is in the group
		newUser := socket.RawMessage{}
		newUser.ParseMessagePack(byte(socket.MessageType.UserJoinedGroup), responseGroup{
			UserID:   user.ID.String(),
			Username: user.Username,
			GroupID:  groupDB.ID.String(),
			IsCPU:    true,
		})

		g.mutex.Lock()
		for i := range g.groups[groupDB.ID] {
			socket.SendQueueMessageSocketID(newUser, g.groups[groupDB.ID][i])
		}
		for k := range g.queue {
			socket.SendQueueMessageSocketID(newUser, k)
		}
		g.mutex.Unlock()
		return
	}

	log.Printf("[Lobby] -> Cannot find a user for the virtual player : %v", user)

}

//AddBot used to add a bot to the group
func (g *groups) AddBot(socketID uuid.UUID) {
	g.mutex.Lock()
	groupID := g.assignment[socketID]
	g.mutex.Unlock()

	//Check if groups exists
	groupDB := model.Group{}
	model.DB().Where("id = ? and status = 0", groupID).First(&groupDB)
	if groupDB.ID != uuid.Nil {
		g.mutex.Lock()
		//Is the group full ?
		if groupDB.PlayersMax-groupDB.VirtualPlayers-len(g.groups[groupID]) > 0 {
			g.mutex.Unlock()
			user := model.User{IsCPU: true}
			model.AddUser(&user)

			user.Username = virtualplayer.AddVirtualPlayer(groupID, user.ID, groupDB.Language)
			log.Printf("[Lobby] -> adding bot in DB: %v", user)
			if user.Username != "" {
				//send response to client
				message := socket.RawMessage{}
				message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
					Response: true,
					Error:    "",
				})

				if socket.SendRawMessageToSocketID(message, socketID) == nil {
					//appel au joueur virtuel
					model.DB().Model(&groupDB).Association("Users").Append(&user)
					groupDB.VirtualPlayers++
					model.DB().Save(&groupDB)
					//Send a message to all the member of the group to advertise that a new user is in the group
					newUser := socket.RawMessage{}
					newUser.ParseMessagePack(byte(socket.MessageType.UserJoinedGroup), responseGroup{
						UserID:   user.ID.String(),
						Username: user.Username,
						GroupID:  groupID.String(),
						IsCPU:    true,
					})

					g.mutex.Lock()
					for i := range g.groups[groupID] {
						socket.SendQueueMessageSocketID(newUser, g.groups[groupID][i])
					}
					for k := range g.queue {
						socket.SendQueueMessageSocketID(newUser, k)
					}
					g.mutex.Unlock()
					return
				}
				return
			}
			g.mutex.Unlock()
			message := socket.RawMessage{}
			message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
				Response: false,
				Error:    "The user could not be found full in virtual player cache.",
			})
			socket.SendQueueMessageSocketID(message, socketID)
			return

		}
		g.mutex.Unlock()
		message := socket.RawMessage{}
		message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
			Response: false,
			Error:    language.MustGetSocket("error.groupIsFull", socketID),
		})
		socket.SendQueueMessageSocketID(message, socketID)
		return

	}
	message := socket.RawMessage{}
	message.ParseMessagePack(byte(socket.MessageType.ResponseJoinGroup), responseGen{
		Response: false,
		Error:    "The group could not be found in DB.",
	})
	socket.SendQueueMessageSocketID(message, socketID)
	return

}

//kickVirtualPlayer quits the groups the virtual player is currently in.
func (g *groups) kickVirtualPlayer(userID uuid.UUID) bool {

	if groupID, username := virtualplayer.KickVirtualPlayer(userID); username != "" && groupID != uuid.Nil {
		message := socket.RawMessage{}
		message.ParseMessagePack(byte(socket.MessageType.ResponseLeaveGroup), responseGroup{
			UserID:   userID.String(),
			Username: username,
			GroupID:  groupID.String(),
			IsCPU:    true,
		})
		g.mutex.Lock()
		for i := range g.groups[groupID] {
			socket.SendQueueMessageSocketID(message, g.groups[groupID][i])
		}
		for k := range g.queue {
			socket.SendQueueMessageSocketID(message, k)
		}
		g.mutex.Unlock()

		return true
	}
	return false
}
