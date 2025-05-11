package db

import "fmt"

func AddVote(roomID, userID, vote string) error {
	_, err := DB.Exec(
		`INSERT INTO votes (room_id, user_id, vote) 
		 VALUES ($1, $2, $3) 
		 ON CONFLICT (room_id, user_id) DO UPDATE SET vote = $3`,
		roomID, userID, vote,
	)
	if err != nil {
		return fmt.Errorf("failed to add vote: %v", err)
	}

	return nil
}

func ResetVotes(roomID string) error {
	_, err := DB.Exec("DELETE FROM votes WHERE room_id = $1", roomID)
	if err != nil {
		return fmt.Errorf("failed to reset votes: %v", err)
	}

	return nil
}

func DeleteVote(roomID, userID string) error {
	_, err := DB.Exec("DELETE FROM votes WHERE room_id = $1 AND user_id = $2", roomID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %v", err)
	}
	return nil
}
