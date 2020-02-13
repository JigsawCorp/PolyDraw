package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

//Connection represents the information of a connection
type Connection struct {
	Base
	UserID        uuid.UUID
	ConnectedAt   int64
	DeconnectedAt int64
}

//MatchPlayed represents the summary of the game
type MatchPlayed struct {
	Base
	UserID        uuid.UUID
	PlayerNames   pq.StringArray `gorm:"type:varchar(20)[]"` // a voir si on laisse le tableau de nom ou si on met un tableau de joueur (moins couteux string)
	MatchDuration int64
	WinnerName    string
}

// Achievement represent an achievement that can be obtained by a player
type Achievement struct {
	Base
	UserID        uuid.UUID
	TropheeName   string
	Description   string
	ObtainingDate int64
}

// UpdateStats updates user's stats each time a match/game is finished
func UpdateStats(userID uuid.UUID, matchPlayed *MatchPlayed) {
	var user User
	DB().Where("userID = ?", userID).First(&user)

	var gamesPlayed int64 = user.GamesPlayed + 1
	var timePlayed int64 = user.TimePlayed + matchPlayed.MatchDuration

	var winRatio float64 = user.WinRatio
	if matchPlayed.WinnerName == user.Username {
		winRatio = (user.WinRatio*float64(gamesPlayed) + 1) / float64(gamesPlayed)
	}

	DB().Where("userID = ?", userID).Updates(map[string]interface{}{
		"gamesPlayed":     gamesPlayed,
		"winRatio":        winRatio,
		"timePlayed":      timePlayed,
		"avgGameDuration": float64(timePlayed) / float64(gamesPlayed),
	})
}