package com.log3900.shared.architecture

enum class EventType {
    // Channel
    ACTIVE_CHANNEL_CHANGED,
    ACTIVE_CHANNEL_MESSAGE_RECEIVED,
    SUBSCRIBED_TO_CHANNEL,
    UNSUBSCRIBED_FROM_CHANNEL,
    CHANNEL_CREATED,
    CHANNEL_DELETED,
    // Message
    RECEIVED_MESSAGE,
    UNREAD_MESSAGES_CHANGED,
    ALL_MESSAGES_CHANGED,
    // Group
    GROUP_CREATED,
    GROUP_DELETED,
    GROUP_UPDATED,
    GROUP_JOINED,
    GROUP_LEFT,
    PLAYER_JOINED_GROUP,
    PLAYER_LEFT_GROUP,
    LEAVE_GROUP,
    // Match
    MATCH_ABOUT_TO_START,
    MATCH_START_RESPONSE,
    MATCH_STARTING,
    MATCH_PLAYERS_UPDATED,
    MATCH_SYNCHRONISATION,
    TURN_TO_DRAW,
    PLAYER_TURN_TO_DRAW,
    PLAYER_GUESSED_WORD,
    GUESSED_WORD_RIGHT,
    GUESSED_WORD_WRONG,
    TIMES_UP,
    TEAMATE_GUESSED_WORD_PROPERLY,
    TEAMATE_GUESSED_WORD_INCORRECTLY,
    MATCH_CANCELLED,
    PLAYER_LEFT_GAME,
    POINTS_GAINED,
    ROUND_ENDED,
    MATCH_ENDED,
    HINT_RESPONSE,
    CHECKPOINT,
    // Session
    LOGOUT,
    SHOW_ERROR_MESSAGE,
    // Settings
    LANGUAGE_CHANGED,
    THEME_CHANGED,
    // MISC
    USER_UPDATED,
    CURRENT_USER_UPDATED
}
data class MessageEvent(var type: EventType, var data: Any?)