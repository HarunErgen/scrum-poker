package room_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/session"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

type JoinRoomRequest struct {
	UserName string `json:"userName"`
}

type JoinRoomResponse struct {
	User *models.User `json:"user"`
	Room *models.Room `json:"room"`
}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	currSession := getSession(r)
	if currSession != nil {
		room, _ := db.GetRoom(currSession.RoomId)
		user, _ := db.GetUser(currSession.UserId)

		resp := JoinRoomResponse{
			User: user,
			Room: room,
		}

		websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
		utils.PrepareJSONResponse(w, http.StatusOK, resp)
		return
	}

	if req.UserName == "" {
		http.Error(w, "User name is required", http.StatusBadRequest)
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	userId := uuid.New().String()
	user := models.NewUser(userId, req.UserName)
	room.AddParticipant(user)

	if err := db.AddParticipantToRoom(roomId, user); err != nil {
		http.Error(w, "Failed to join room", http.StatusInternalServerError)
		return
	}

	sessionID, _ := session.GlobalManager.CreateSession(userId, roomId)

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	resp := JoinRoomResponse{
		User: user,
		Room: room,
	}
	websocket.GlobalHub.BroadcastRoomUpdate(roomId, room)
	utils.PrepareJSONResponse(w, http.StatusOK, resp)
}

func getSession(r *http.Request) *models.Session {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		fmt.Print("Session cookie not found.")
		return nil
	}
	sessionId := cookie.Value
	currSession, err := db.GetSession(sessionId)
	if err != nil {
		fmt.Print("Session not found in database.")
		return nil
	}

	if currSession.IsExpired() {
		fmt.Print("Session is expired.")
		err := db.DeleteSession(sessionId)
		if err != nil {
			return nil
		}
		return nil
	}

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	if roomId != currSession.RoomId {
		fmt.Print("Room Id does not match session.")
		return nil
	}

	if err := db.UpdateUserOnlineStatus(currSession.UserId, true); err != nil {
		fmt.Print("Failed to update user status.")
		return nil
	}

	currSession.Refresh(session.SessionTTL)

	if err := db.UpdateSession(currSession); err != nil {
		fmt.Print("Failed to update session.")
		return nil
	}
	return currSession
}
