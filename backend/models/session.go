package models

import (
	"time"
)

type Session struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	RoomId    string    `json:"roomId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewSession(id, userId, roomId string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		Id:        id,
		UserId:    userId,
		RoomId:    roomId,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Refresh(ttl time.Duration) {
	s.ExpiresAt = time.Now().Add(ttl)
}

func (s *Session) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":        s.Id,
		"userId":    s.UserId,
		"roomId":    s.RoomId,
		"createdAt": s.CreatedAt,
		"expiresAt": s.ExpiresAt,
	}
}
