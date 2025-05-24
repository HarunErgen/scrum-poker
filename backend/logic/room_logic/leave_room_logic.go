package room_logic

import (
	"fmt"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

func LeaveRoom(roomId, userId string, broadcastFunc models.BroadcastFunc) error {

	room, err := db.GetRoom(roomId)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if _, ok := room.Participants[userId]; !ok {
		return fmt.Errorf("user not in room")
	}

	if room.ScrumMaster == userId {
		participantsCopy := make(map[string]*models.User)
		for id, user := range room.Participants {
			participantsCopy[id] = user
		}
		delete(participantsCopy, userId)

		if len(participantsCopy) > 0 {
			room.AssignRandomScrumMaster(participantsCopy)

			if err := db.UpdateScrumMaster(roomId, room.ScrumMaster); err != nil {
				return fmt.Errorf("failed to update Scrum Master: %w", err)
			}

			message := &models.Message{
				Action: models.ActionTypeTransfer,
				Payload: map[string]interface{}{
					"userId":           userId,
					"newScrumMasterId": room.ScrumMaster,
				},
			}
			broadcastFunc(roomId, message)
		}
	}

	room.RemoveParticipant(userId)

	if err := db.RemoveParticipantFromRoom(roomId, userId); err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}

	if err := db.DeleteUser(userId); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if len(room.Participants) == 0 {
		if err := db.DeleteRoom(roomId); err != nil {
			return fmt.Errorf("failed to delete room: %w", err)
		}
	}
	return nil
}
