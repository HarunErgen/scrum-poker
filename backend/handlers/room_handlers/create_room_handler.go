package room_handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"net/http"
	"time"
)

type CreateRoomRequest struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
}

type RoomResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	CreatedAt     time.Time              `json:"created_at"`
	ScrumMaster   string                 `json:"scrum_master"`
	Participants  map[string]interface{} `json:"participants"`
	Votes         map[string]string      `json:"votes"`
	VotesRevealed bool                   `json:"votes_revealed"`
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}
	if req.UserName == "" {
		http.Error(w, "User name is required", http.StatusBadRequest)
		return
	}

	roomID := uuid.New().String()
	userID := uuid.New().String()

	user := models.NewUser(userID, req.UserName)
	room := models.NewRoom(roomID, req.Name, userID)
	room.AddParticipant(user)

	if err := db.CreateRoom(room); err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	if err := db.AddParticipantToRoom(roomID, user); err != nil {
		http.Error(w, "Failed to add participant to room", http.StatusInternalServerError)
		return
	}

	resp := RoomResponse{
		ID:            room.ID,
		Name:          room.Name,
		CreatedAt:     room.CreatedAt,
		ScrumMaster:   room.ScrumMaster,
		Participants:  map[string]interface{}{userID: user.ToJSON()},
		Votes:         make(map[string]string),
		VotesRevealed: false,
	}
	utils.PrepareJSONResponse(w, http.StatusCreated, resp)
}
