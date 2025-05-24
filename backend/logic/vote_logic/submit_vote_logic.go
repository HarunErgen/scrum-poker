package vote_logic

import (
	"fmt"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

func SubmitVote(userId, roomId, vote string) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if _, ok := room.Participants[userId]; !ok {
		return fmt.Errorf("user %s not in room %s", userId, roomId)
	}

	if vote == "" {
		if err := db.DeleteVote(roomId, userId); err != nil {
			return fmt.Errorf("failed to delete vote: %w", err)
		}
	} else {
		if !models.IsValidVote(vote) {
			return fmt.Errorf("invalid vote value: %s", vote)
		}

		if err := db.AddVote(roomId, userId, vote); err != nil {
			return fmt.Errorf("failed to submit vote: %w", err)
		}
	}
	return nil
}
