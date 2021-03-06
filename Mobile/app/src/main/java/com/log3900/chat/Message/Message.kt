package com.log3900.chat.Message

import com.squareup.moshi.Json
import java.util.*

data class ReceivedMessage(@Json(name = "Message") var message: String, @Json(name = "ChannelID") var channelID:  UUID,
                           @Json(name = "UserID") var userID: UUID, @Json(name = "Username") var username: String,
                           @Json(name = "Timestamp") var timestamp: Date)

data class SentMessage(@Json(name = "Message")var message: String, @Json(name = "ChannelID") var channelID: UUID)

class EventMessage(var type: Type, var message: String) {
    enum class Type {
        USER_LEFT_CHANNEL,
        USER_JOINED_CHANNEL,
        USERNAME_CHANGED
    }
}