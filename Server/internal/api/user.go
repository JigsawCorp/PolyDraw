package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/jigsawcorp/log3900/internal/context"
	"gitlab.com/jigsawcorp/log3900/internal/language"
	"gitlab.com/jigsawcorp/log3900/internal/services/messenger"
	"gitlab.com/jigsawcorp/log3900/internal/socket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/jigsawcorp/log3900/model"
	"gitlab.com/jigsawcorp/log3900/pkg/rbody"
	"gitlab.com/jigsawcorp/log3900/pkg/secureb"
)

type singleUserResponse struct {
	ID        string
	FirstName string
	LastName  string
	Username  string
	Email     string
	PictureID int
	CreatedAt int64
	UpdatedAt int64
	IsCPU     bool
}

type socketUserChange struct {
	UserID    string
	NewName   string
	PictureID int
	IsCPU     bool
	OldName   string
}

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	model.DB().Find(&users)
	json.NewEncoder(w).Encode(users)
}

//GetUser return a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user model.User
	model.DB().Where("id = ?", vars["id"]).First(&user)
	if user.ID != uuid.Nil {
		json.NewEncoder(w).Encode(singleUserResponse{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			PictureID: user.PictureID,
			IsCPU:     user.IsCPU,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.CreatedAt.Unix(),
		})
	} else {
		rbody.JSONError(w, http.StatusNotFound, language.MustGetRest("error.userNotFound", r))
	}
}

//PutUser update the user
func PutUser(w http.ResponseWriter, r *http.Request) {
	var request authRegisterRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		rbody.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	userid := r.Context().Value(context.CtxUserID)
	var user model.User
	var oldUsername string
	model.DB().Where("id = ?", userid).First(&user)
	if user.ID != uuid.Nil {
		updated := false
		if request.Email != "" {
			//Validate email
			if !regexEmail.MatchString(request.Email) {
				rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.emailInvalid", r))
				return
			}
			user.Email = request.Email
			updated = true
		}

		if request.Username != "" {
			//Validate username
			username := strings.ToLower(request.Username)
			if !regexUsername.MatchString(username) {
				rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.usernameInvalid", r))
				return
			}

			var count int64
			model.DB().Model(&model.User{}).Where("username = ?", username).Count(&count)
			if count > 0 {
				rbody.JSONError(w, http.StatusConflict, language.MustGetRest("error.usernameExists", r))
				return
			}
			oldUsername = user.Username
			user.Username = request.Username
			updated = true
		}

		if request.FirstName != "" {
			user.FirstName = request.FirstName
			updated = true
		}

		if request.LastName != "" {
			user.LastName = request.LastName
			updated = true
		}

		if request.Password != "" {
			if len(request.Password) < 8 {
				rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.passwordInvalid", r))
				return
			}

			hash, err := secureb.HashPassword(request.Password)
			if err != nil {
				rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.userUpdate", r))
				return
			}
			user.HashedPassword = hash
			updated = true
		}

		if request.PictureID != 0 {
			if request.PictureID < 1 || request.PictureID > 16 {
				rbody.JSONError(w, http.StatusBadRequest, "Invalid picture id, the number must be between 1 and 16.")
				return
			}
			user.PictureID = request.PictureID
			updated = true
		}

		if updated {
			model.DB().Save(&user)
			json.NewEncoder(w).Encode("ok")

			if request.Username != "" || request.PictureID != 0 {
				//Broadcast to all users
				message := socket.RawMessage{}
				message.ParseMessagePack(byte(socket.MessageType.UsernameChange), socketUserChange{
					UserID:    user.ID.String(),
					PictureID: user.PictureID,
					OldName:   oldUsername,
					NewName:   user.Username,
					IsCPU:     user.IsCPU,
				})
				messenger.BroadcastAll(message)
			}
		} else {
			rbody.JSONError(w, http.StatusBadRequest, language.MustGetRest("error.modifications", r))
		}

	} else {
		rbody.JSONError(w, http.StatusNotFound, language.MustGetRest("error.userNotFound", r))
	}
}
