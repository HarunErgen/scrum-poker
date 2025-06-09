package session_handlers

import (
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/config"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/session"
	"github.com/scrum-poker/backend/utils"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	roomId := vars["roomId"]

	if userId == "" || roomId == "" {
		http.Error(w, "User Id and Room Id are required", http.StatusBadRequest)
		return
	}

	_, err := db.GetUser(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if _, ok := room.Participants[userId]; !ok {
		http.Error(w, "User not in room", http.StatusForbidden)
		return
	}

	sessionID, err := session.GlobalManager.CreateSession(userId, roomId)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.Cfg.Cookie.Secure,
		SameSite: config.Cfg.Cookie.SameSite,
	})

	if err := db.UpdateUserOnlineStatus(userId, true); err != nil {
		http.Error(w, "Failed to update user status", http.StatusInternalServerError)
		return
	}

	utils.PrepareJSONResponse(w, http.StatusCreated, map[string]string{
		"sessionId": sessionID,
	})
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	sessionId := cookie.Value
	if sessionId == "" {
		http.Error(w, "Session id is required", http.StatusBadRequest)
		return
	}

	existingSession, err := db.GetSession(sessionId)
	if err != nil {
		http.Error(w, "Session not found in database", http.StatusNotFound)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	if existingSession.RoomId != roomId {
		http.Error(w, "Room id does not match session", http.StatusForbidden)
		return
	}

	if existingSession.IsExpired() {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(existingSession.UserId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	room, err := db.GetRoom(existingSession.RoomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	existingSession.Refresh(session.TTL)
	if err := db.UpdateSession(existingSession); err != nil {
		http.Error(w, "Failed to update session", http.StatusInternalServerError)
		return
	}

	if err := db.UpdateUserOnlineStatus(existingSession.UserId, true); err != nil {
		http.Error(w, "Failed to update user status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"session": existingSession.ToJSON(),
		"user":    user.ToJSON(),
		"room":    room.ToJSON(),
	}
	message := &models.Message{
		Action: models.ActionTypeOnline,
		Payload: map[string]interface{}{
			"userId": existingSession.UserId,
		},
	}
	websocket.GlobalHub.Broadcast(roomId, message)
	utils.PrepareJSONResponse(w, http.StatusOK, response)
}

func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	sessionID := cookie.Value
	if sessionID == "" {
		http.Error(w, "Session Id is required", http.StatusBadRequest)
		return
	}

	if err := session.GlobalManager.DeleteSession(sessionID); err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   config.Cfg.Cookie.Secure,
		SameSite: config.Cfg.Cookie.SameSite,
	})

	utils.PrepareJSONResponse(w, http.StatusOK, []byte("OK"))
}
