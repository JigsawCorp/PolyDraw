package model

import (
	"github.com/google/uuid"
)

//GroupStatus
//0 created and waiting for players
//1 game started
//2 game ended
//3 abandoned for example 5 minutes of inactivity or server closing

//Group represent a new group used for a game
type Group struct {
	Base
	Name           string
	Owner          User `gorm:"foreignkey:OwnerID"`
	OwnerID        uuid.UUID
	Language       int `gorm:"default:0"` // 0 -> EN | 1 -> FR
	PlayersMax     int
	GameType       int //0 -> FFA | 1 -> solo | 2 -> coop
	VirtualPlayers int
	NbRound        int
	Difficulty     int     `gorm:"default:0"`
	Status         int     `gorm:"default:0"`
	Users          []*User `gorm:"many2many:group_membership"`
}
