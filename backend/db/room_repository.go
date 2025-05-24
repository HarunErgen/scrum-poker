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
		room.Id, room.Name, room.CreatedAt, room.ScrumMaster,
	)
	if err != nil {
		return fmt.Errorf("failed to create room: %v", err)
	}

	return nil
}

func DeleteRoom(roomId string) error {
	_, err := DB.Exec("DELETE FROM rooms WHERE id = $1", roomId)
	if err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}
	return nil
}

func GetRoom(roomId string) (*models.Room, error) {
	var room models.Room

	err := DB.QueryRow(
		"SELECT id, name, created_at, scrum_master FROM rooms WHERE id = $1",
		roomId,
	).Scan(&room.Id, &room.Name, &room.CreatedAt, &room.ScrumMaster)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	room.Participants = make(map[string]*models.User)
	room.Votes = make(map[string]string)
	room.VotesRevealed = false

	rows, err := DB.Query(`
		SELECT u.id, u.name, u.created_at, u.is_online
		FROM users u
		JOIN room_participants rp ON u.id = rp.user_id
		WHERE rp.room_id = $1
	`, roomId)
	if err != nil {
		return nil, fmt.Errorf("failed to get room participants: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := new(models.User)
		var userCreatedAt time.Time
		err := rows.Scan(&user.Id, &user.Name, &userCreatedAt, &user.IsOnline)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		user.CreatedAt = userCreatedAt
		room.Participants[user.Id] = user
	}

	voteRows, err := DB.Query("SELECT user_id, vote FROM votes WHERE room_id = $1", roomId)
	if err != nil {
		return nil, fmt.Errorf("failed to get room votes: %v", err)
	}
	defer voteRows.Close()

	for voteRows.Next() {
		var userId, vote string
		err := voteRows.Scan(&userId, &vote)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %v", err)
		}
		room.Votes[userId] = vote
	}

	return &room, nil
}

func AddParticipantToRoom(roomId string, user *models.User) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.Id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}

	if count == 0 {
		_, err = tx.Exec(
			"INSERT INTO users (id, name, created_at, is_online) VALUES ($1, $2, $3, $4)",
			user.Id, user.Name, user.CreatedAt, user.IsOnline,
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	} else {
		_, err = tx.Exec(
			"UPDATE users SET name = $1, is_online = $2 WHERE id = $3",
			user.Name, user.IsOnline, user.Id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO room_participants (room_id, user_id) 
		 VALUES ($1, $2) 
		 ON CONFLICT (room_id, user_id) DO NOTHING`,
		roomId, user.Id,
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

func RemoveParticipantFromRoom(roomId, userId string) error {
	_, err := DB.Exec(
		"DELETE FROM room_participants WHERE room_id = $1 AND user_id = $2",
		roomId, userId,
	)
	if err != nil {
		return fmt.Errorf("failed to remove participant from room: %v", err)
	}

	_, err = DB.Exec(
		"DELETE FROM votes WHERE room_id = $1 AND user_id = $2",
		roomId, userId,
	)
	if err != nil {
		return fmt.Errorf("failed to remove votes: %v", err)
	}

	return nil
}

func UpdateScrumMaster(roomId, newScrumMasterID string) error {
	_, err := DB.Exec(
		"UPDATE rooms SET scrum_master = $1 WHERE id = $2",
		newScrumMasterID, roomId,
	)
	if err != nil {
		return fmt.Errorf("failed to update Scrum Master: %v", err)
	}
	return nil
}

func GetAllRooms() ([]*models.Room, error) {
	rows, err := DB.Query("SELECT id, name, created_at, scrum_master FROM rooms")
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %v", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		var room models.Room
		var createdAt time.Time
		err := rows.Scan(&room.Id, &room.Name, &createdAt, &room.ScrumMaster)
		if err != nil {
			return nil, fmt.Errorf("failed to scan room: %v", err)
		}
		room.CreatedAt = createdAt
		room.Participants = make(map[string]*models.User)
		room.Votes = make(map[string]string)
		room.VotesRevealed = false

		participantRows, err := DB.Query(`
			SELECT u.id, u.name, u.created_at, u.is_online
			FROM users u
			JOIN room_participants rp ON u.id = rp.user_id
			WHERE rp.room_id = $1
		`, room.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get room participants: %v", err)
		}
		defer participantRows.Close()

		for participantRows.Next() {
			user := new(models.User)
			var userCreatedAt time.Time
			err := participantRows.Scan(&user.Id, &user.Name, &userCreatedAt, &user.IsOnline)
			if err != nil {
				return nil, fmt.Errorf("failed to scan user: %v", err)
			}
			user.CreatedAt = userCreatedAt
			room.Participants[user.Id] = user
		}

		voteRows, err := DB.Query("SELECT user_id, vote FROM votes WHERE room_id = $1", room.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get room votes: %v", err)
		}
		defer voteRows.Close()

		for voteRows.Next() {
			var userId, vote string
			err := voteRows.Scan(&userId, &vote)
			if err != nil {
				return nil, fmt.Errorf("failed to scan vote: %v", err)
			}
			room.Votes[userId] = vote
		}

		rooms = append(rooms, &room)
	}

	return rooms, nil
}

func GetRoomByUserId(userId string) (*models.Room, error) {
	var room models.Room

	err := DB.QueryRow(
		`SELECT r.id, r.name, r.created_at, r.scrum_master
				FROM rooms r 
				JOIN room_participants rp ON r.id = rp.room_id 
				WHERE rp.user_id = $1`,
		userId,
	).Scan(&room.Id, &room.Name, &room.CreatedAt, &room.ScrumMaster)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	room.Participants = make(map[string]*models.User)
	room.Votes = make(map[string]string)
	room.VotesRevealed = false

	rows, err := DB.Query(`
		SELECT u.id, u.name, u.created_at, u.is_online
		FROM users u
		JOIN room_participants rp ON u.id = rp.user_id
		WHERE rp.room_id = $1
	`, room.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get room participants: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := new(models.User)
		var userCreatedAt time.Time
		err := rows.Scan(&user.Id, &user.Name, &userCreatedAt, &user.IsOnline)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		user.CreatedAt = userCreatedAt
		room.Participants[user.Id] = user
	}

	voteRows, err := DB.Query("SELECT user_id, vote FROM votes WHERE room_id = $1", room.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get room votes: %v", err)
	}
	defer voteRows.Close()

	for voteRows.Next() {
		var userId, vote string
		err := voteRows.Scan(&userId, &vote)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %v", err)
		}
		room.Votes[userId] = vote
	}

	return &room, nil
}
