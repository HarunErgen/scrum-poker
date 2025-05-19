package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/scrum-poker/backend/models"
)

func DeleteUser(userId string) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func GetUser(userId string) (*models.User, error) {
	var user models.User
	var createdAt time.Time
	err := DB.QueryRow(
		"SELECT id, name, created_at, is_online FROM users WHERE id = $1",
		userId,
	).Scan(&user.Id, &user.Name, &createdAt, &user.IsOnline)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	user.CreatedAt = createdAt
	return &user, nil
}

func UpdateUserOnlineStatus(userId string, isOnline bool) error {
	_, err := DB.Exec(
		"UPDATE users SET is_online = $1 WHERE id = $2",
		isOnline, userId,
	)
	if err != nil {
		return fmt.Errorf("failed to update user online status: %v", err)
	}
	return nil
}

func UpdateUser(user *models.User) error {
	_, err := DB.Exec(
		"UPDATE users SET name = $1, is_online = $2 WHERE id = $3",
		user.Name, user.IsOnline, user.Id,
	)
	if err != nil {
		return err
	}
	return nil
}
