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
	UserId string `json:"userId"`
}

func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	var req LeaveRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserId == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}
	userId := req.UserId

	room, err := db.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if _, ok := room.Participants[userId]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	if room.ScrumMaster == userId {
		participantsCopy := make(map[string]*models.User)
		for id, user := range room.Participants {
			participantsCopy[id] = user
		}
		delete(participantsCopy, userId)

		if len(participantsCopy) > 0 {
			room.AssignRandomScrumMaster(participantsCopy)

			if err := db.UpdateScrumMaster(roomId, room.ScrumMaster); err != nil {
				http.Error(w, "Failed to update Scrum Master", http.StatusInternalServerError)
				return
			}
		}
	}

	room.RemoveParticipant(userId)

	if err := db.RemoveParticipantFromRoom(roomId, userId); err != nil {
		http.Error(w, "Failed to leave room", http.StatusInternalServerError)
		return
	}

	if err := db.DeleteUser(userId); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if len(room.Participants) == 0 {
		if err := db.DeleteRoom(roomId); err != nil {
			http.Error(w, "Failed to delete room", http.StatusInternalServerError)
			return
		}
	}

	websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
