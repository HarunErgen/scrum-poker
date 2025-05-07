package vote_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type SubmitVoteRequest struct {
	UserID string `json:"user_id"`
	Vote   string `json:"vote"`
}

func SubmitVoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var req SubmitVoteRequest
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

	if _, ok := room.Participants[req.UserID]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	if req.Vote == "" {
		room.RemoveVote(req.UserID)

		if err := db.DeleteVote(roomID, req.UserID); err != nil {
			http.Error(w, "Failed to delete vote", http.StatusInternalServerError)
			return
		}

		websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
		utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
		return
	}

	if !models.IsValidVote(req.Vote) {
		http.Error(w, "Invalid vote", http.StatusBadRequest)
		return
	}

	room.AddVote(req.UserID, req.Vote)

	if err := db.AddVote(roomID, req.UserID, req.Vote); err != nil {
		http.Error(w, "Failed to submit vote", http.StatusInternalServerError)
		return
	}

	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
