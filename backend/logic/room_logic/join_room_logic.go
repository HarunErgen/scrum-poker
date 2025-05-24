package room_logic

import (
	"github.com/google/uuid"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type DatabaseError struct {
	Operation string
	Message   string
}

func (e DatabaseError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Resource string
	Message  string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func JoinRoom(roomId string, userName string, existingSession *models.Session, broadcastFunc models.BroadcastFunc) (string, error) {
	if existingSession != nil {
		user, err := db.GetUser(existingSession.UserId)
		if err != nil {
			return "", DatabaseError{
				Operation: "GetUser",
				Message:   "Failed to get user information",
			}
		}

		message := &models.Message{
			Action:  models.ActionTypeJoin,
			Payload: user.ToJSON(),
		}
		broadcastFunc(roomId, message)

		return existingSession.UserId, nil
	}

	if userName == "" {
		return "", ValidationError{
			Field:   "userName",
			Message: "User name is required",
		}
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		return "", NotFoundError{
			Resource: "Room",
			Message:  "Room not found",
		}
	}

	userId := uuid.New().String()
	user := models.NewUser(userId, userName)
	room.AddParticipant(user)

	if err := db.AddParticipantToRoom(roomId, user); err != nil {
		return "", DatabaseError{
			Operation: "AddParticipantToRoom",
			Message:   "Failed to join room",
		}
	}

	message := &models.Message{
		Action:  models.ActionTypeJoin,
		Payload: user.ToJSON(),
	}
	broadcastFunc(roomId, message)

	return userId, nil
}
