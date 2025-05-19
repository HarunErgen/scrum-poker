package handlers

import (
	"fmt"
	"github.com/scrum-poker/backend/handlers/room_handlers"
	"github.com/scrum-poker/backend/handlers/session_handlers"
	"github.com/scrum-poker/backend/handlers/user_handlers"
	"github.com/scrum-poker/backend/handlers/vote_handlers"
	"github.com/scrum-poker/backend/handlers/websocket_handlers"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

var (
	UpdateUserHandler = user_handlers.UpdateUserHandler
)

var (
	CreateRoomHandler          = room_handlers.CreateRoomHandler
	GetRoomHandler             = room_handlers.GetRoomHandler
	JoinRoomHandler            = room_handlers.JoinRoomHandler
	LeaveRoomHandler           = room_handlers.LeaveRoomHandler
	TransferScrumMasterHandler = room_handlers.TransferScrumMasterHandler
)

var (
	SubmitVoteHandler  = vote_handlers.SubmitVoteHandler
	RevealVotesHandler = vote_handlers.RevealVotesHandler
	ResetVotesHandler  = vote_handlers.ResetVotesHandler
)

var (
	WebSocketHandler = websocket_handlers.WebSocketHandler
)

var (
	CreateSessionHandler = session_handlers.CreateSessionHandler
	GetSessionHandler    = session_handlers.GetSessionHandler
	DeleteSessionHandler = session_handlers.DeleteSessionHandler
)
