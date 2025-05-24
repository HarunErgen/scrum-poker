package room_logic

import (
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

func GetRoom(roomId string) (*models.Room, error) {
	return db.GetRoom(roomId)
}
