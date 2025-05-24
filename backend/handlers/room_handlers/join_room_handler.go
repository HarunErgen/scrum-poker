package room_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/logic/room_logic"
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
	UserId string `json:"userId"`
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

	userId, err := room_logic.JoinRoom(roomId, req.UserName, currSession, websocket.GlobalHub.Broadcast)
	if err != nil {
		switch e := err.(type) {
		case room_logic.ValidationError:
			http.Error(w, e.Error(), http.StatusBadRequest)
		case room_logic.NotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		case room_logic.DatabaseError:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		return
	}

	if currSession == nil {
		sessionID, _ := session.GlobalManager.CreateSession(userId, roomId)

		http.SetCookie(w, &http.Cookie{
			Name:     "sessionId",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
	}

	utils.PrepareJSONResponse(w, http.StatusOK, JoinRoomResponse{
		UserId: userId,
	})
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

	currSession.Refresh(session.TTL)

	if err := db.UpdateSession(currSession); err != nil {
		fmt.Print("Failed to update session.")
		return nil
	}
	return currSession
}
