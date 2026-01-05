package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/session"
)

var GlobalHub *Hub

type Hub struct {
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
}

func Init() {
	GlobalHub = &Hub{
		rooms: make(map[string]map[*Client]bool),
	}

	session.InitSessionManager(
		GlobalHub.Broadcast,
		GlobalHub.IsUserConnected,
	)
}

func (h *Hub) RegisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, exists := h.rooms[c.roomId]
	if !exists {
		clients = make(map[*Client]bool)
		h.rooms[c.roomId] = clients
	}

	for existing := range clients {
		if existing.userId == c.userId {
			delete(clients, existing)
			close(existing.send)
			if existing.conn != nil {
				existing.conn.Close()
			}
		}
	}

	clients[c] = true

	go h.notifyUserOnline(c.roomId, c.userId)
}

func (h *Hub) UnregisterClient(c *Client) {
	h.mu.Lock()

	clients, exists := h.rooms[c.roomId]
	if !exists {
		h.mu.Unlock()
		return
	}

	if _, ok := clients[c]; !ok {
		h.mu.Unlock()
		return
	}

	delete(clients, c)
	close(c.send)

	userStillConnected := false
	for existing := range clients {
		if existing.userId == c.userId {
			userStillConnected = true
			break
		}
	}

	if len(clients) == 0 {
		delete(h.rooms, c.roomId)
	}

	h.mu.Unlock()

	if !userStillConnected {
		go h.handleUserOffline(c.roomId, c.userId)
	}
}

func (h *Hub) notifyUserOnline(roomId, userId string) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		log.Printf("Error getting room %s: %v", roomId, err)
		return
	}

	if _, ok := room.Participants[userId]; !ok {
		return
	}

	h.Broadcast(roomId, &models.Message{
		Action:  models.ActionTypeOnline,
		Payload: map[string]interface{}{"userId": userId},
	})
}

func (h *Hub) handleUserOffline(roomId, userId string) {
	time.Sleep(100 * time.Millisecond)

	if h.IsUserConnected(roomId, userId) {
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		log.Printf("Error getting room %s: %v", roomId, err)
		return
	}

	if _, ok := room.Participants[userId]; !ok {
		return
	}

	// Refresh or create session for potential reconnection
	existingSession, err := db.GetSessionByUserID(userId)
	if err == nil && existingSession != nil {
		existingSession.Refresh(session.TTL)
		if err := db.UpdateSession(existingSession); err != nil {
			log.Printf("Error updating session: %v", err)
		}
	} else {
		if _, err := session.GlobalManager.CreateSession(userId, roomId); err != nil {
			log.Printf("Error creating session: %v", err)
		}
	}

	h.Broadcast(roomId, &models.Message{
		Action:  models.ActionTypeOffline,
		Payload: map[string]interface{}{"userId": userId},
	})
}

func (h *Hub) IsUserConnected(roomId, userId string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, exists := h.rooms[roomId]
	if !exists {
		return false
	}

	for c := range clients {
		if c.userId == userId {
			return true
		}
	}
	return false
}

func (h *Hub) Broadcast(roomId string, msg *models.Message) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	h.mu.RLock()
	clients, exists := h.rooms[roomId]
	if !exists {
		h.mu.RUnlock()
		return
	}

	clientList := make([]*Client, 0, len(clients))
	for c := range clients {
		clientList = append(clientList, c)
	}
	h.mu.RUnlock()

	for _, c := range clientList {
		select {
		case c.send <- msgBytes:
		default:
			go h.UnregisterClient(c)
		}
	}
}

// GetConnectedUserIds returns all connected user IDs for a room
func (h *Hub) GetConnectedUserIds(roomId string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, exists := h.rooms[roomId]
	if !exists {
		return nil
	}

	userIds := make([]string, 0, len(clients))
	for c := range clients {
		userIds = append(userIds, c.userId)
	}
	return userIds
}
