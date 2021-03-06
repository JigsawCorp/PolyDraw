package com.log3900.profile.stats

import com.google.gson.JsonObject
import com.log3900.profile.ProfileRestService
import com.log3900.settings.language.LanguageManager
import com.log3900.user.account.AccountRepository
import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory

// Service?
/**
 * Stores all statistics. Getters should be called from a coroutine
 */
object StatsRepository {
    enum class StatsCategory {
        ALL,
        GENERAL,
        CONNECTION,
        MATCHES_PLAYED,
        ACHIEVEMENTS
    }

//    private var subscribers
    private lateinit var userStats: UserStats
    private lateinit var historyStats: HistoryStats

    /**
     * Custom getters for the stats.
     * Does a Rest call if not initialized.
     */
    private suspend fun getUserStats(): UserStats {
//        if (!StatsRepository::userStats.isInitialized) {
//            fetchUserStats()
//        }
        fetchUserStats()
        return userStats
    }

    private suspend fun getHistoryStats(): HistoryStats {
//        if (!StatsRepository::historyStats.isInitialized) {
//            fetchHistoryStats()
//        }
        fetchHistoryStats()
        return historyStats
    }

    suspend fun getAllUserStats(): UserStats = getUserStats()
    suspend fun getGamesPlayedHistory(): List<GamePlayed> = getHistoryStats().gamesPlayedHistory ?: listOf()
    suspend fun getConnectionHistory(): List<Connection> = getHistoryStats().connectionHistory
    suspend fun getAchievements(): List<Achievement> = listOf() //getHistoryStats().achievements

    private suspend fun fetchUserStats() {
        userStats = sendUserStatsRequest()
    }

    private suspend fun fetchHistoryStats() {
        historyStats = sendHistoryStatsRequest()
    }

    private suspend fun sendUserStatsRequest(): UserStats {
        val session = AccountRepository.getInstance().getAccount().sessionToken
        val language = LanguageManager.getCurrentLanguageCode()
        val responseJson = ProfileRestService.service.getUserStats(session, language)   //TODO: get language

        if (responseJson.isSuccessful && responseJson.body() != null) {
            val json = responseJson.body()!!
            return parseJsonToStats(json)
        } else {
            throw Exception("${responseJson.code()} : ${responseJson.errorBody()?.string()}")
        }
    }

    private suspend fun sendHistoryStatsRequest(): HistoryStats {
        val session = AccountRepository.getInstance().getAccount().sessionToken
        val language = LanguageManager.getCurrentLanguageCode()
        val responseJson = ProfileRestService.service.getHistory(session, language)

        if (responseJson.isSuccessful && responseJson.body() != null) {
            val json = responseJson.body()!!
            return parseJsonToStats(json)
        } else {
            throw Exception("${responseJson.code()} : ${responseJson.errorBody()?.string()}")
        }
    }

    private inline fun <reified StatsType> parseJsonToStats(json: JsonObject): StatsType {
        println(json.toString())
        val moshi = Moshi.Builder()
            .add(KotlinJsonAdapterFactory())
            .build()
        val adapter: JsonAdapter<StatsType> = moshi.adapter(StatsType::class.java)

        return adapter.fromJson(json.toString())!!
    }
}