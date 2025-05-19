package vote_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type ResetVotesRequest struct {
	UserId string `json:"userId"`
}

func ResetVotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	var req ResetVotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserId == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.ScrumMaster != req.UserId {
		http.Error(w, "Only the Scrum Master can reset votes", http.StatusForbidden)
		return
	}

	room.ResetVotes()

	if err := db.ResetVotes(roomId); err != nil {
		http.Error(w, "Failed to reset votes", http.StatusInternalServerError)
		return
	}

	websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
