package com.log3900.profile.stats

import com.squareup.moshi.Json

// TODO: use Dates?
// TODO: Best score solo
data class UserStats(
    @Json(name = "GamesPlayed")             val gamesPlayed:            Int,
    @Json(name = "WinRatio")                val winRatio:              Double,
    @Json(name = "AvgGameDuration")         val averageGameDuration:    Double,
    @Json(name = "TimePlayed")              val timePlayed:             Long,
    @Json(name = "BestScoreSolo")           val bestScoreSolo:          Int
)

data class HistoryStats(
    @Json(name = "MatchesPlayedHistory")    val gamesPlayedHistory:     List<GamePlayed>?,
    @Json(name = "ConnectionHistory")       val connectionHistory:      List<Connection>
//    @Json(name = "Achievements")            val achievements:           List<Achievement>
)

data class GamePlayed(
    @Json(name = "MatchDuration")   val duration:       Int,
    @Json(name = "WinnerName")      val winner:         String,
    @Json(name = "MatchType")       val matchType:      String
//    @Json(name = "PlayersNames")    val playerNames:    List<PlayerName>
)

data class Connection(
    @Json(name = "ConnectedAt")     val connectedAt:    Int,
    @Json(name = "DisconnectedAt")  val disconnectedAt: Int
)

data class Achievement(
    @Json(name = "TropheeName")     val title:          String,
    @Json(name = "Description")     val description:    String,
    @Json(name = "ObtainingDate")   val unlockDate:     Int
)

data class PlayerName(
    @Json(name = "PlayerName") val name: String
)
