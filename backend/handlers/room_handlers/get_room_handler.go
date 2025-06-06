package room_handlers

import (
	"github.com/gorilla/mux"
	"github.com/scrum-poker/backend/logic/room_logic"
	"github.com/scrum-poker/backend/utils"
	"net/http"
)

func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	room, err := room_logic.GetRoom(roomId)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	utils.PrepareJSONResponse(w, http.StatusOK, room.ToJSON())
}
