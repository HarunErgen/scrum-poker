package room_handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type JoinRoomRequest struct {
	UserName string `json:"user_name"`
}

type JoinRoomResponse struct {
	User *models.User `json:"user"`
	Room *models.Room `json:"room"`
}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserName == "" {
		http.Error(w, "User name is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	userID := uuid.New().String()
	user := models.NewUser(userID, req.UserName)
	room.AddParticipant(user)

	if err := db.AddParticipantToRoom(roomID, user); err != nil {
		http.Error(w, "Failed to join room", http.StatusInternalServerError)
		return
	}

	resp := JoinRoomResponse{
		User: user,
		Room: room,
	}
	websocket.CommonHub.BroadcastRoomUpdate(roomID, room)
	utils.PrepareJSONResponse(w, http.StatusOK, resp)
}
