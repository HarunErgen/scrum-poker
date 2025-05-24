package room_logic

import (
	"errors"
	"github.com/google/uuid"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

func CreateRoom(roomName, userName string) (*models.Room, *models.User, error) {
	if roomName == "" {
		return nil, nil, errors.New("room name is required")
	}
	if userName == "" {
		return nil, nil, errors.New("user name is required")
	}

	roomId := uuid.New().String()
	userId := uuid.New().String()

	user := models.NewUser(userId, userName)
	room := models.NewRoom(roomId, roomName, userId)
	room.AddParticipant(user)

	if err := db.CreateRoom(room); err != nil {
		return nil, nil, errors.New("failed to create room")
	}

	if err := db.AddParticipantToRoom(roomId, user); err != nil {
		return nil, nil, errors.New("failed to add participant to room")
	}

	return room, user, nil
}
