package vote_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type RevealVotesRequest struct {
	UserID string `json:"user_id"`
}

func RevealVotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var req RevealVotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.ScrumMaster != req.UserID {
		http.Error(w, "Only the Scrum Master can reveal votes", http.StatusForbidden)
		return
	}

	room.RevealVotes()

	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
