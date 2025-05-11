package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/scrum-poker/backend/models"
)

func CreateRoom(room *models.Room) error {
	_, err := DB.Exec(
		"INSERT INTO rooms (id, name, created_at, scrum_master) VALUES ($1, $2, $3, $4)",
		room.ID, room.Name, room.CreatedAt, room.ScrumMaster,
	)
	if err != nil {
		return fmt.Errorf("failed to create room: %v", err)
	}

	return nil
}

func DeleteRoom(roomID string) error {
	_, err := DB.Exec("DELETE FROM rooms WHERE id = $1", roomID)
	if err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}
	return nil
}

func GetRoom(roomID string) (*models.Room, error) {
	var room models.Room
	var createdAt time.Time

	err := DB.QueryRow(
		"SELECT id, name, created_at, scrum_master FROM rooms WHERE id = $1",
		roomID,
	).Scan(&room.ID, &room.Name, &createdAt, &room.ScrumMaster)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	room.CreatedAt = createdAt
	room.Participants = make(map[string]*models.User)
	room.Votes = make(map[string]string)
	room.VotesRevealed = false

	rows, err := DB.Query(`
		SELECT u.id, u.name, u.created_at
		FROM users u
		JOIN room_participants rp ON u.id = rp.user_id
		WHERE rp.room_id = $1
	`, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room participants: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := new(models.User)
		var userCreatedAt time.Time
		err := rows.Scan(&user.ID, &user.Name, &userCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		user.CreatedAt = userCreatedAt
		room.Participants[user.ID] = user
	}

	voteRows, err := DB.Query("SELECT user_id, vote FROM votes WHERE room_id = $1", roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room votes: %v", err)
	}
	defer voteRows.Close()

	for voteRows.Next() {
		var userID, vote string
		err := voteRows.Scan(&userID, &vote)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %v", err)
		}
		room.Votes[userID] = vote
	}

	return &room, nil
}

func AddParticipantToRoom(roomID string, user *models.User) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}

	if count == 0 {
		_, err = tx.Exec(
			"INSERT INTO users (id, name, created_at) VALUES ($1, $2, $3)",
			user.ID, user.Name, user.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	} else {
		_, err = tx.Exec(
			"UPDATE users SET name = $1 WHERE id = $2",
			user.Name, user.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO room_participants (room_id, user_id) 
		 VALUES ($1, $2) 
		 ON CONFLICT (room_id, user_id) DO NOTHING`,
		roomID, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to add participant to room: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func RemoveParticipantFromRoom(roomID, userID string) error {
	_, err := DB.Exec(
		"DELETE FROM room_participants WHERE room_id = $1 AND user_id = $2",
		roomID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove participant from room: %v", err)
	}

	_, err = DB.Exec(
		"DELETE FROM votes WHERE room_id = $1 AND user_id = $2",
		roomID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove votes: %v", err)
	}

	return nil
}

func TransferScrumMaster(roomID, newScrumMasterID string) error {
	_, err := DB.Exec(
		"UPDATE rooms SET scrum_master = $1 WHERE id = $2",
		newScrumMasterID, roomID,
	)
	if err != nil {
		return fmt.Errorf("failed to transfer Scrum Master: %v", err)
	}
	return nil
}

func UpdateScrumMaster(roomID, newScrumMasterID string) error {
	_, err := DB.Exec(
		"UPDATE rooms SET scrum_master = $1 WHERE id = $2",
		newScrumMasterID, roomID,
	)
	if err != nil {
		return fmt.Errorf("failed to update Scrum Master: %v", err)
	}
	return nil
}
