package room_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
	"time"
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

	if room.ScrumMaster == userID {
		participantsCopy := make(map[string]*models.User)
		for id, user := range room.Participants {
			participantsCopy[id] = user
		}
		delete(participantsCopy, userID)

		if len(participantsCopy) > 0 {
			newScrumMasterID := selectRandomNextScrumMaster(participantsCopy)
			room.ScrumMaster = newScrumMasterID

			if err := db.UpdateScrumMaster(roomID, newScrumMasterID); err != nil {
				http.Error(w, "Failed to update Scrum Master", http.StatusInternalServerError)
				return
			}
		}
	}

	if _, ok := room.Participants[userID]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	room.RemoveParticipant(userID)

	if err := db.RemoveParticipantFromRoom(roomID, userID); err != nil {
		http.Error(w, "Failed to leave room", http.StatusInternalServerError)
		return
	}

	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}

func selectRandomNextScrumMaster(participants map[string]*models.User) string {
	var candidates []string
	for userID := range participants {
		candidates = append(candidates, userID)
	}

	if len(candidates) == 0 {
		return ""
	}

	randomIndex := time.Now().UnixNano() % int64(len(candidates))
	return candidates[randomIndex]
}
