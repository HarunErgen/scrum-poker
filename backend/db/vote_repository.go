package db

import "fmt"

func AddVote(roomId, userId, vote string) error {
	_, err := DB.Exec(
		`INSERT INTO votes (room_id, user_id, vote) 
		 VALUES ($1, $2, $3) 
		 ON CONFLICT (room_id, user_id) DO UPDATE SET vote = $3`,
		roomId, userId, vote,
	)
	if err != nil {
		return fmt.Errorf("failed to add vote: %v", err)
	}

	return nil
}

func ResetVotes(roomId string) error {
	_, err := DB.Exec("DELETE FROM votes WHERE room_id = $1", roomId)
	if err != nil {
		return fmt.Errorf("failed to reset votes: %v", err)
	}

	return nil
}

func DeleteVote(roomId, userId string) error {
	_, err := DB.Exec("DELETE FROM votes WHERE room_id = $1 AND user_id = $2", roomId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %v", err)
	}
	return nil
}
