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
	UserId           string `json:"userId"`
	NewScrumMasterID string `json:"newScrumMasterId"`
}

func TransferScrumMasterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	var req TransferScrumMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserId == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}
	if req.NewScrumMasterID == "" {
		http.Error(w, "New Scrum Master Id is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.ScrumMaster != req.UserId {
		http.Error(w, "Only the current Scrum Master can transfer the role", http.StatusForbidden)
		return
	}

	if _, ok := room.Participants[req.NewScrumMasterID]; !ok {
		http.Error(w, "New Scrum Master is not in the room", http.StatusBadRequest)
		return
	}

	room.TransferScrumMaster(req.NewScrumMasterID)

	if err := db.TransferScrumMaster(roomId, req.NewScrumMasterID); err != nil {
		http.Error(w, "Failed to transfer Scrum Master", http.StatusInternalServerError)
		return
	}

	websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
