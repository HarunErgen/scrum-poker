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
	UserName string `json:"userName"`
}

type RoomResponse struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	CreatedAt     time.Time              `json:"createdAt"`
	ScrumMaster   string                 `json:"scrumMaster"`
	Participants  map[string]interface{} `json:"participants"`
	Votes         map[string]string      `json:"votes"`
	VotesRevealed bool                   `json:"votesRevealed"`
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

	roomId := uuid.New().String()
	userId := uuid.New().String()

	user := models.NewUser(userId, req.UserName)
	room := models.NewRoom(roomId, req.Name, userId)
	room.AddParticipant(user)

	if err := db.CreateRoom(room); err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	if err := db.AddParticipantToRoom(roomId, user); err != nil {
		http.Error(w, "Failed to add participant to room", http.StatusInternalServerError)
		return
	}

	resp := RoomResponse{
		Id:            room.Id,
		Name:          room.Name,
		CreatedAt:     room.CreatedAt,
		ScrumMaster:   room.ScrumMaster,
		Participants:  map[string]interface{}{userId: user.ToJSON()},
		Votes:         make(map[string]string),
		VotesRevealed: false,
	}
	utils.PrepareJSONResponse(w, http.StatusCreated, resp)
}
