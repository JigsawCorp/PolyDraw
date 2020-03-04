package router

import (
	"gitlab.com/jigsawcorp/log3900/internal/services/drawing"
	"gitlab.com/jigsawcorp/log3900/internal/services/healthcheck"
	"gitlab.com/jigsawcorp/log3900/internal/services/lobby"
	"gitlab.com/jigsawcorp/log3900/internal/services/match"
	"gitlab.com/jigsawcorp/log3900/internal/services/messenger"
	"gitlab.com/jigsawcorp/log3900/internal/socket"
)

func (r *Router) routing() {
	//Messenger
	r.handle(socket.MessageType.MessageSent, messenger.BMessageSent)
	r.handle(socket.MessageType.CreateChannel, messenger.BCreateChannel)
	r.handle(socket.MessageType.DestroyChannel, messenger.BDestroyChannel)
	r.handle(socket.MessageType.JoinChannel, messenger.BJoinChannel)
	r.handle(socket.MessageType.LeaveChannel, messenger.BLeaveChannel)

	//Healthcheck
	r.handle(socket.MessageType.HealthCheckResponse, healthcheck.BReceived)

	//Drawing
	r.handle(socket.MessageType.PreviewDrawing, drawing.BPreview)

	//Lobby
	r.handle(socket.MessageType.RequestLeaveGroup, lobby.BLeaveGroup)
	r.handle(socket.MessageType.RequestJoinGroup, lobby.BJoinGroup)

	//Match
	r.handle(socket.MessageType.RequestReadyMatch, match.BMatchReady)
	r.handle(socket.MessageType.RequestQuitMatch, match.BMatchQuit)
	r.handle(socket.MessageType.RequestGuessWordMatch, match.BMatchGuess)
	r.handle(socket.MessageType.RequestHintMatch, match.BMatchHint)
}
