package com.log3900.socket

enum class Event(var eventType: Int) {
    SOCKET_CONNECTION(0),
    SERVER_RESPONSE(1),
    CLIENT_DISCONNECT(2),
    DISCONNEECTED_BY_SERVER(3),
    HEALTH_CHECK_SERVER(9),
    HEALTH_CHECK_CLIENT(10),
    MESSAGE_SENT(20),
    MESSAGE_RECEIVED(21),
    JOIN_CHANNEL(22),
    JOINED_CHANNEL(23),
    LEAVE_CHANNEL(24),
    LEFT_CHANNEL(25),
    CREATE_CHANNEL(26),
    CHANNEL_CREATED(27),
    DELETE_CHANNEL(28),
    CHANNEL_DELETED(29),
    STROKE_DATA_CLIENT(30),
    STROKE_DATA_SERVER(31),
    DRAW_START_CLIENT(32),
    DRAW_START_SERVER(33),
    DRAW_END_CLIENT(34),
    DRAW_END_SERVER(35),
    DRAW_PREVIEW_REQUEST(36),   // TODO: Remove, test purposes only
    DRAW_PREVIEW_RESPONSE(37),  // TODO: Remove, test purposes only
    STROKE_ERASE_CLIENT(38),
    STROKE_ERASE_SERVER(39),
    JOIN_GROUP_REQUEST(40),
    JOIN_GROUP_RESPONSE(41),
    USER_JOINED_GROUP(43),
    LEAVE_GROUP_REQUEST(44),
    USER_LEFT_GROUP(45),
    START_MATCH(48),
    START_MATCH_RESPONSE(49),
    GROUP_CREATED(51),
    GROUP_DELETED(53),
    KICK_PLAYER(54),
    ADD_VIRTUAL_PLAYER(56),
    ADD_VIRTUAL_PLAYER_RESPONSE(57),
    MATCH_ABOUT_TO_START(61),
    READY_TO_PLAY_MATCH(62),
    MATCH_STARTING(63),
    LEAVE_MATCH(64),
    PLAYER_LEFT_MATCH(65),
    PLAYER_TURN_TO_DRAW(67),
    TURN_TO_DRAW(69),
    TIMES_UP(71),
    PLAYER_SYNC(73),
    GUESS_WORD(74),
    GUESS_WORD_RESPONSE(75),
    PLAYER_GUESSED_WORD(77),
    CHECKPOINT(79),
    MATCH_ENDED(81),
    HINT_REQUEST(82),
    HINT_RESPONSE(83),
    ROUND_ENDED(85),
    TEAMATE_GUESSED_WORD_PROPERLY(87),
    TEAMATE_GUESSED_WORD_INCORRECTLY(89),
    MATCH_CANCELLED(91),
    // MISC
    ACHIEVEMENT_RECEIVED(101),
    USERNAME_CHANGED(110),
    LANGUAGE_CHANGED(112),
    SERVER_ERROR(255)
}

data class Message(var type: Event, var data: ByteArray)