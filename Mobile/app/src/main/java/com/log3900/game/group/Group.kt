package com.log3900.game.group

import com.log3900.R
import com.squareup.moshi.Json
import java.util.*

enum class Difficulty {
    EASY,
    MEDIUM,
    HARD;

    companion object {
        fun stringRes(item: Difficulty): Int {
            when (item) {
                EASY -> return R.string.difficulty_easy_title
                MEDIUM -> return R.string.difficulty_medium_title
                HARD -> return R.string.difficulty_hard_title
            }
        }

        fun timeToPlay(item: Difficulty): Long {
            when (item) {
                EASY -> return 1000 * 60
                MEDIUM -> return 1000 * 45
                HARD -> return 1000 * 30
            }
        }
    }
}

enum class MatchMode {
    FFA,
    SOLO,
    COOP;

    companion object {
        fun stringRes(item: MatchMode): Int {
            when (item) {
                FFA -> return R.string.match_mode_ffa_title
                COOP -> return R.string.match_mode_coop_title
                SOLO -> return R.string.match_mode_solo_title
            }
        }
        fun imageRes(item: MatchMode): Int {
            when (item) {
                FFA -> return R.drawable.ic_swords
                SOLO -> return R.drawable.ic_solo_sprint_logo
                COOP -> return R.drawable.ic_treasure
            }
        }
    }
}

enum class Language(val languageCode: String) {
    ENGLISH("EN"),
    FRENCH("FR");

    companion object {
        fun stringRes(item: Language): Int {
            when (item) {
                FRENCH -> return R.string.french
                ENGLISH -> return R.string.english
            }
        }
    }

}

class Group(@Json(name = "ID") var ID: UUID, @Json(name = "GroupName") var groupName: String,
            @Json(name = "PlayersMax") var playersMax: Int, @Json(name = "NbRound") var rounds: Int?,
            @Json(name = "GameType") var gameType: MatchMode, @Json(name = "Difficulty") var difficulty: Difficulty,
            @Json(name = "OwnerID") var ownerID: UUID, @Json(name = "OwnerName") var ownerName: String,
            @Json(name = "Language") var language: Language, @Json(name = "Players") var players: ArrayList<Player>) {


    companion object {
        fun minAmountOfPlayers(matchMode: MatchMode): Int {
            when (matchMode) {
                MatchMode.SOLO -> return 1
                MatchMode.COOP -> return 2
                MatchMode.FFA -> return 2
            }
        }

        fun maxAmountOfPlayers(matchMode: MatchMode): Int {
            when (matchMode) {
                MatchMode.SOLO -> return 1
                MatchMode.COOP -> return 4
                MatchMode.FFA -> return 5
            }
        }
    }

    fun getDuration(): Long {
        if (gameType == MatchMode.FFA) {
            return (1000 * 60 * players.size * rounds!!).toLong()
        } else {
            return Difficulty.timeToPlay(difficulty)
        }
    }
}

class UserLeftGroup(
    var userID: UUID,
    var username: String,
    var groupID: UUID,
    var isCPU: Boolean,
    var isKicked: Boolean
)