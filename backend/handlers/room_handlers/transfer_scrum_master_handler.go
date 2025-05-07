package room_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type TransferScrumMasterRequest struct {
	UserID           string `json:"user_id"`
	NewScrumMasterID string `json:"new_scrum_master_id"`
}

func TransferScrumMasterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var req TransferScrumMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if req.NewScrumMasterID == "" {
		http.Error(w, "New Scrum Master ID is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.ScrumMaster != req.UserID {
		http.Error(w, "Only the current Scrum Master can transfer the role", http.StatusForbidden)
		return
	}

	if _, ok := room.Participants[req.NewScrumMasterID]; !ok {
		http.Error(w, "New Scrum Master is not in the room", http.StatusBadRequest)
		return
	}

	room.TransferScrumMaster(req.NewScrumMasterID)

	if err := db.TransferScrumMaster(roomID, req.NewScrumMasterID); err != nil {
		http.Error(w, "Failed to transfer Scrum Master", http.StatusInternalServerError)
		return
	}

	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
