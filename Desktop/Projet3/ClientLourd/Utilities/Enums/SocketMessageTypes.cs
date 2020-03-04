﻿namespace ClientLourd.Utilities.Enums
{
    public enum SocketMessageTypes
    {
        ServerConnection = 0,
        ServerConnectionResponse = 1,
        ServerDisconnection = 2,
        UserDisconnected = 3,
        HealthCheck = 9,
        HealthCheckResponse = 10,
        MessageSent = 20,
        MessageReceived = 21,
        JoinChannel = 22,
        UserJoinedChannel = 23,
        LeaveChannel = 24,
        UserLeftChannel = 25,
        CreateChannel = 26,
        UserCreatedChannel = 27,
        DeleteChannel = 28,
        UserDeletedChannel = 29,
        UserStrokeSent = 30,
        ServerStrokeSent = 31,
        StartDrawing = 32,
        ServerStartsDrawing = 33,
        EndDrawing = 34,
        ServerEndsDrawing = 35,
        DrawingPreviewRequest = 36,
        DrawingPreviewResponse = 37,
        DeleteStroke = 38, 
        UserDeletedStroke = 39,
        ServerJoinLobby = 41,
        LobbyCreated = 51,
        ServerMessage = 255,
    }
}