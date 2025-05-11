package room_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type LeaveRoomRequest struct {
	UserID string `json:"user_id"`
}

func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var req LeaveRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	room, err := db.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if _, ok := room.Participants[userID]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	if room.ScrumMaster == userID {
		participantsCopy := make(map[string]*models.User)
		for id, user := range room.Participants {
			participantsCopy[id] = user
		}
		delete(participantsCopy, userID)

		if len(participantsCopy) > 0 {
			room.AssignRandomScrumMaster(participantsCopy)

			if err := db.UpdateScrumMaster(roomID, room.ScrumMaster); err != nil {
				http.Error(w, "Failed to update Scrum Master", http.StatusInternalServerError)
				return
			}
		}
	}

	room.RemoveParticipant(userID)

	if err := db.RemoveParticipantFromRoom(roomID, userID); err != nil {
		http.Error(w, "Failed to leave room", http.StatusInternalServerError)
		return
	}

	if err := db.DeleteUser(userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if len(room.Participants) == 0 {
		if err := db.DeleteRoom(roomID); err != nil {
			http.Error(w, "Failed to delete room", http.StatusInternalServerError)
			return
		}
	}

	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
