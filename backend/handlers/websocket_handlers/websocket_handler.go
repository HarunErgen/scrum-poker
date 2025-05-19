package websocket_handlers

import (
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/websocket"
	"net/http"
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
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

	websocket.ServeWs(websocket.GlobalHub, w, r, roomId, userId)
}
