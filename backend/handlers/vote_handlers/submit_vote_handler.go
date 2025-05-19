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
	UserId string `json:"userId"`
	Vote   string `json:"vote"`
}

func SubmitVoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	var req SubmitVoteRequest
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

	if _, ok := room.Participants[req.UserId]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	if req.Vote == "" {
		room.RemoveVote(req.UserId)

		if err := db.DeleteVote(roomId, req.UserId); err != nil {
			http.Error(w, "Failed to delete vote", http.StatusInternalServerError)
			return
		}

		websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
		utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
		return
	}

	if !models.IsValidVote(req.Vote) {
		http.Error(w, "Invalid vote", http.StatusBadRequest)
		return
	}

	room.AddVote(req.UserId, req.Vote)

	if err := db.AddVote(roomId, req.UserId, req.Vote); err != nil {
		http.Error(w, "Failed to submit vote", http.StatusInternalServerError)
		return
	}

	websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
