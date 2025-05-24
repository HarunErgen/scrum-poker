package room_logic

import (
	"fmt"
	"github.com/scrum-poker/backend/db"
)

func TransferScrumMaster(userId, roomId, newScrumMasterId string) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if room.ScrumMaster != userId {
		return fmt.Errorf("only the Scrum Master can transfer the role")
	}

	if _, ok := room.Participants[newScrumMasterId]; !ok {
		return fmt.Errorf("new Scrum Master is not in the room")
	}

	if err := db.UpdateScrumMaster(roomId, newScrumMasterId); err != nil {
		return fmt.Errorf("failed to transfer Scrum Master: %w", err)
	}
	return nil
}
