package api

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jigsawcorp/log3900/internal/context"
	"gitlab.com/jigsawcorp/log3900/internal/language"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moby/moby/pkg/namesgenerator"
	"gitlab.com/jigsawcorp/log3900/internal/services/lobby"
	"gitlab.com/jigsawcorp/log3900/model"
	"gitlab.com/jigsawcorp/log3900/pkg/rbody"
)

type requestGroupCreate struct {
	GroupName  string
	PlayersMax int
	NbRound    int
	GameType   int
	Difficulty int
}

type responsePlayer struct {
	ID       string
	Username string
	IsCPU    bool
}

type responsePlayerReal struct {
	ID       string
	Username string
}

type responseGroupCreate struct {
	GroupName string
	GroupID   string
}
type responseGroup struct {
	ID         string
	GroupName  string
	PlayersMax int
	GameType   int
	Difficulty int
	Status     int
	OwnerName  string
	OwnerID    string
	Language   int
	NbRound    int
	Players    []responsePlayer
}

const maxPlayer = 8

//PostGroup used to create a new group
func PostGroup(w http.ResponseWriter, r *http.Request) {
	var request requestGroupCreate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		rbody.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if request.Difficulty < 0 || request.Difficulty > 2 {
		rbody.JSONError(w, http.StatusBadRequest, "The difficulty must be between 0 and 2.")
		return
	}

	if request.GameType > 2 || request.GameType < 0 {
		rbody.JSONError(w, http.StatusBadRequest, "The game mode must be between 0 and 2")
		return
	}

	if (request.PlayersMax > maxPlayer || request.PlayersMax < 1) && request.GameType != 1 {
		rbody.JSONError(w, http.StatusBadRequest, fmt.Sprintf(language.MustGetRest("error.playerCount", r), maxPlayer))
		return
	}
	if request.GameType == 0 && (request.NbRound <= 0 || request.NbRound > 5) {
		rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.ffaRound", r))
		return
	}
	if request.PlayersMax != 1 && request.GameType == 1 {
		rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.soloCount", r))
		return
	}

	if request.PlayersMax > maxPlayer {
		rbody.JSONError(w, http.StatusBadRequest, fmt.Sprintf(language.MustGetRest("error.playerMax", r), maxPlayer))
		return
	}
	userid := r.Context().Value(context.CtxUserID).(uuid.UUID)

	var groupName string
	if request.GroupName != "" {
		groupName = request.GroupName
	} else {
		//Generate a docker like name
		groupName = namesgenerator.GetRandomName(0)
	}

	group := model.Group{
		OwnerID:    userid,
		Name:       groupName,
		NbRound:    request.NbRound,
		PlayersMax: request.PlayersMax,
		GameType:   request.GameType,
		Difficulty: request.Difficulty,
		Language:   r.Context().Value(context.CtxLang).(int),
		Status:     0,
	}
	var user model.User
	model.DB().Model(&user).Where("id = ?", userid).First(&user)
	model.DB().Create(&group)

	group.Owner = user
	response := responseGroupCreate{
		GroupName: groupName,
		GroupID:   group.ID.String(),
	}
	lobby.Instance().CreateGroup(&group)

	rbody.JSON(w, http.StatusOK, response)
}

//GetGroups returns all the groups that are currently available
func GetGroups(w http.ResponseWriter, r *http.Request) {
	var groups []model.Group
	model.DB().Model(&groups).Related(&model.User{}, "Users")
	model.DB().Preload("Users").Preload("Owner").Where("status = ?", 0).Find(&groups)

	response := make([]responseGroup, 0, len(groups))
	for i := range groups {
		lenUsers := len(groups[i].Users)
		if lenUsers > 0 {
			players := make([]responsePlayer, lenUsers)
			for j := range groups[i].Users {
				players[j] = responsePlayer{
					ID:       groups[i].Users[j].ID.String(),
					Username: groups[i].Users[j].Username,
					IsCPU:    groups[i].Users[j].IsCPU,
				}
			}

			response = append(response, responseGroup{
				ID:         groups[i].ID.String(),
				GroupName:  groups[i].Name,
				PlayersMax: groups[i].PlayersMax,
				GameType:   groups[i].GameType,
				Difficulty: groups[i].Difficulty,
				Status:     0,
				Language:   groups[i].Language,
				NbRound:    groups[i].NbRound,
				OwnerName:  groups[i].Owner.Username,
				OwnerID:    groups[i].OwnerID.String(),
				Players:    players,
			})
		}
	}
	rbody.JSON(w, http.StatusOK, response)
}

//GetGroup returns the details of the specific group
func GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID, err := uuid.Parse(vars["id"])
	if err != nil {
		rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.invalidUUID", r))
		return
	}
	var group model.Group
	model.DB().Model(&group).Related(&model.User{}, "Users")
	model.DB().Preload("Users").Preload("Owner").Where("id = ?", groupID).First(&group)

	if group.ID != uuid.Nil {
		players := make([]responsePlayer, len(group.Users))
		for j := range group.Users {
			players[j] = responsePlayer{
				ID:       group.Users[j].ID.String(),
				Username: group.Users[j].Username,
				IsCPU:    group.Users[j].IsCPU,
			}
		}

		response := responseGroup{
			ID:         group.ID.String(),
			GroupName:  group.Name,
			PlayersMax: group.PlayersMax,
			GameType:   group.GameType,
			Difficulty: group.Difficulty,
			Status:     group.Status,
			Language:   group.Language,
			NbRound:    group.NbRound,
			OwnerName:  group.Owner.Username,
			OwnerID:    group.Owner.ID.String(),
			Players:    players,
		}
		rbody.JSON(w, http.StatusOK, response)
		return
	}

	rbody.JSONError(w, http.StatusNotFound, language.MustGetRest("error.groupNotFound", r))
	return

}
