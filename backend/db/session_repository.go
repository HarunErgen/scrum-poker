package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/scrum-poker/backend/models"
)

func CreateSession(session *models.Session) error {
	_, err := DB.Exec(
		"INSERT INTO sessions (id, user_id, room_id, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)",
		session.Id, session.UserId, session.RoomId, session.CreatedAt, session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	return nil
}

func GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	err := DB.QueryRow(
		"SELECT id, user_id, room_id, created_at, expires_at FROM sessions WHERE id = $1",
		sessionID,
	).Scan(&session.Id, &session.UserId, &session.RoomId, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return &session, nil
}

func UpdateSession(session *models.Session) error {
	_, err := DB.Exec(
		"UPDATE sessions SET expires_at = $1 WHERE id = $2",
		session.ExpiresAt, session.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}
	return nil
}

func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}

func GetSessionsByRoomID(roomId string) ([]*models.Session, error) {
	rows, err := DB.Query(
		"SELECT id, user_id, room_id, created_at, expires_at FROM sessions WHERE room_id = $1",
		roomId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions for room: %v", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := new(models.Session)
		err := rows.Scan(&session.Id, &session.UserId, &session.RoomId, &session.CreatedAt, &session.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func GetSessionByUserID(userId string) (*models.Session, error) {
	var session models.Session
	err := DB.QueryRow(
		"SELECT id, user_id, room_id, created_at, expires_at FROM sessions WHERE user_id = $1",
		userId,
	).Scan(&session.Id, &session.UserId, &session.RoomId, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return &session, nil
}
