package models

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	IsOnline  bool      `json:"isOnline"`
}

func NewUser(id, name string) *User {
	return &User{
		Id:        id,
		Name:      name,
		CreatedAt: time.Now(),
		IsOnline:  true,
	}
}

func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.Id,
		"name":      u.Name,
		"createdAt": u.CreatedAt,
		"isOnline":  u.IsOnline,
	}
}
