package room_handlers

import (
	"encoding/json"
	"github.com/scrum-poker/backend/logic/room_logic"
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

	room, user, err := room_logic.CreateRoom(req.Name, req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := RoomResponse{
		Id:            room.Id,
		Name:          room.Name,
		CreatedAt:     room.CreatedAt,
		ScrumMaster:   room.ScrumMaster,
		Participants:  map[string]interface{}{user.Id: user.ToJSON()},
		Votes:         make(map[string]string),
		VotesRevealed: false,
	}
	utils.PrepareJSONResponse(w, http.StatusCreated, resp)
}
