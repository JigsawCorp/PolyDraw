package model

import "github.com/google/uuid"

//Game represent a game
type Game struct {
	Base
	Word       string
	Difficulty int
	Hints      []*GameHint `gorm:"foreignkey:GameID"`
	Image      *GameImage  //Represent the hash of the file
	Language   int
}

//GameHint represents a game hint
type GameHint struct {
	GameID uuid.UUID `gorm:"primary_key;auto_increment:false"`
	Hint   string    `gorm:"primary_key;auto_increment:false"`
}

//GameImage represents the game image
type GameImage struct {
	Game         Game      `gorm:"foreignkey:GameID"`
	GameID       uuid.UUID `gorm:"primary_key"`
	Mode         int
	BrushSize    int
	BlackLevel   float64
	OriginalFile string
	SVGFile      string
	ImageFile    string
}
