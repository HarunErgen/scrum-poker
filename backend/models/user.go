package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

func NewUser(id, name string) *User {
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		IsActive:  true,
	}
}

func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"created_at": u.CreatedAt,
		"is_active":  u.IsActive,
	}
}
