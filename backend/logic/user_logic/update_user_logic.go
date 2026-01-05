package user_logic

import (
	"fmt"

	"github.com/scrum-poker/backend/db"
)

func RenameUser(userId, roomId, newName string) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if _, ok := room.Participants[userId]; !ok {
		return fmt.Errorf("user not in room")
	}

	if err := db.UpdateUserName(userId, newName); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
