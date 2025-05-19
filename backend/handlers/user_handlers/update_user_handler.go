package user_handlers

import (
	"encoding/json"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Id == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}

	existingUser, err := db.GetUser(user.Id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	existingUser.Name = user.Name
	existingUser.IsOnline = user.IsOnline

	err = db.UpdateUser(existingUser)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	room, err := db.GetRoomByUserId(existingUser.Id)
	if err != nil {
		http.Error(w, "Failed to get room by user id: "+err.Error(), http.StatusInternalServerError)
		return
	}
	websocket.GlobalHub.BroadcastRoomUpdate(room.Id, room)
	utils.PrepareJSONResponse(w, http.StatusOK, existingUser.ToJSON())
}
