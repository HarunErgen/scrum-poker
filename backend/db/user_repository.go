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
		"SELECT id, name, created_at FROM users WHERE id = $1",
		userId,
	).Scan(&user.Id, &user.Name, &createdAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	user.CreatedAt = createdAt
	return &user, nil
}

func UpdateUserName(userId, name string) error {
	_, err := DB.Exec(
		"UPDATE users SET name = $1 WHERE id = $2",
		name, userId,
	)
	if err != nil {
		return fmt.Errorf("failed to update user name: %v", err)
	}
	return nil
}
