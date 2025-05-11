package db

import (
	"fmt"
)

func DeleteUser(userID string) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
