package session

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
)

const (
	SessionTTL      = 3 * time.Minute
	CleanupInterval = 1 * time.Minute
)

type BroadcastFunc func(roomId string, room *models.Room)

var GlobalManager *Manager

type Manager struct {
	broadcastFunc BroadcastFunc
}

func NewManager(broadcastFunc BroadcastFunc) *Manager {
	return &Manager{
		broadcastFunc: broadcastFunc,
	}
}

func InitSessionManager(broadcastFunc BroadcastFunc) {
	GlobalManager = NewManager(broadcastFunc)
	GlobalManager.StartCleanupProcess()
	log.Println("Session manager initialized and cleanup process started")
}

func (m *Manager) CreateSession(userId, roomId string) (string, error) {
	sessionID := uuid.New().String()
	session := models.NewSession(sessionID, userId, roomId, SessionTTL)

	if err := db.CreateSession(session); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (m *Manager) DeleteSession(sessionID string) error {
	session, err := db.GetSession(sessionID)
	if err != nil {
		return err
	}

	if err := db.UpdateUserOnlineStatus(session.UserId, false); err != nil {
		return err
	}

	room, err := db.GetRoom(session.RoomId)
	if err == nil {
		m.broadcastFunc(session.RoomId, room)
	}

	return db.DeleteSession(sessionID)
}

func (m *Manager) StartCleanupProcess() {
	ticker := time.NewTicker(CleanupInterval)
	go func() {
		for {
			<-ticker.C
			m.cleanupExpiredSessions()
		}
	}()
}

func (m *Manager) cleanupExpiredSessions() {
	rooms, err := db.GetAllRooms()
	if err != nil {
		log.Printf("Error getting rooms: %v", err)
		return
	}

	for _, room := range rooms {
		sessions, err := db.GetSessionsByRoomID(room.Id)
		if err != nil {
			log.Printf("Error getting sessions for room %s: %v", room.Id, err)
			continue
		}

		for _, session := range sessions {
			if session.IsExpired() {
				userId := session.UserId

				if user, exists := room.Participants[userId]; exists && user.IsOnline {
					log.Printf("User %s is online, refreshing session %s", userId, session.Id)
					session.Refresh(SessionTTL)
					if err := db.UpdateSession(session); err != nil {
						log.Printf("Error refreshing session: %v", err)
					}
					continue
				} else if exists {
					log.Printf("Session %s expired, cleaning up", session.Id)

					if err := db.RemoveParticipantFromRoom(room.Id, userId); err != nil {
						log.Printf("Error removing participant from room: %v", err)
					}

					if room.ScrumMaster == userId && len(room.Participants) > 0 {
						participantsCopy := make(map[string]*models.User)
						for id, user := range room.Participants {
							if id != userId {
								participantsCopy[id] = user
							}
						}

						if len(participantsCopy) > 0 {
							room.AssignRandomScrumMaster(participantsCopy)
							if err := db.UpdateScrumMaster(room.Id, room.ScrumMaster); err != nil {
								log.Printf("Error updating scrum master: %v", err)
							}
						}
					}

					if err := db.DeleteUser(userId); err != nil {
						log.Printf("Error deleting user: %v", err)
					}
					if err := db.DeleteSession(session.Id); err != nil {
						log.Printf("Error deleting session: %v", err)
					}

					updatedRoom, err := db.GetRoom(room.Id)
					if err == nil {
						m.broadcastFunc(room.Id, updatedRoom)
					}
				}
			}
		}

		if len(room.Participants) == 0 {
			log.Printf("Room %s is empty, deleting", room.Id)
			if err := db.DeleteRoom(room.Id); err != nil {
				log.Printf("Error deleting room: %v", err)
			}
		}
	}
}
