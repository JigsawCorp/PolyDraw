package stats

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gitlab.com/jigsawcorp/log3900/internal/match"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"gitlab.com/jigsawcorp/log3900/model"
)

// DataStats stats of user
type DataStats struct {
	GamesPlayed     int64
	WinRatio        float64
	AvgGameDuration float64
	TimePlayed      int64
	BestScoreSolo   int64
}

func updateMatchesPlayed(stats match.StatsData) {
	model.DB().Create(stats.Match)
	for _, socketID := range stats.SocketsID {
		userID, err := auth.GetUserID(socketID)
		if err != nil {
			log.Printf("[Stats] -> [Error] Can't find userID from socketID: %v.", socketID)
			continue
		}
		model.DB().Create(&model.MatchPlayedMembership{UserID: userID, MatchID: stats.Match.ID})
	}
}

func setDisconnection(userID uuid.UUID) {
	var c model.Connection
	model.DB().Model(&model.Connection{}).Where("user_id = ?", userID).Order("created_at desc").Offset(0).Limit(1).Find(&c)
	if c.ID != uuid.Nil {
		model.DB().Model(&model.Connection{}).Where("id = ?", c.ID).Update("disconnected_at", time.Now().Unix())
	}
}

func createConnection(userID uuid.UUID) {
	model.DB().Create(&model.Connection{UserID: userID, ConnectedAt: time.Now().Unix()})
}

// GetStats find in BD all stats of user
func GetStats(userID uuid.UUID) (DataStats, string) {
	var user model.User
	model.DB().Model(&model.User{}).Where("id = ?", userID).Find(&user)

	if user.ID == uuid.Nil {
		return DataStats{}, "No user found with this userID"
	}

	var matches []model.MatchPlayed
	model.DB().Model(&model.MatchPlayed{}).Joins("JOIN match_played_memberships ON match_played_memberships.match_id = match_playeds.id AND match_played_memberships.user_id = ?", userID).Find(&matches)

	if len(matches) == 0 {
		return DataStats{AvgGameDuration: 0, GamesPlayed: 0, TimePlayed: 0, WinRatio: 0, BestScoreSolo: 0}, ""
	}

	gamesPlayed := int64(len(matches))
	nbWins := 0
	nbFFA := 0
	timePlayed := int64(0)
	bestScoreSolo := int64(0)
	for _, match := range matches {
		//Cherche minimum score solo
		if match.MatchType == 1 && match.PointsSoloCoop > bestScoreSolo {
			bestScoreSolo = match.PointsSoloCoop
		}
		if match.MatchType == 0 {
			nbFFA++
			if user.ID.String() == match.WinnerID {
				nbWins++
			}
		}
		timePlayed += match.MatchDuration
	}

	winRatio := float64(0)

	if nbFFA != 0 {
		winRatio = float64(nbWins) / float64(nbFFA)
	}

	return DataStats{AvgGameDuration: float64(timePlayed) / float64(gamesPlayed), GamesPlayed: gamesPlayed, TimePlayed: timePlayed, WinRatio: winRatio, BestScoreSolo: bestScoreSolo}, ""
}
