package vote_logic

import (
	"fmt"
	"github.com/scrum-poker/backend/db"
)

func ResetVotes(userId, roomId string) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if room.ScrumMaster != userId {
		return fmt.Errorf("only the Scrum Master can reset votes")
	}

	if err := db.ResetVotes(roomId); err != nil {
		return fmt.Errorf("failed to reset votes: %w", err)
	}
	return nil
}
