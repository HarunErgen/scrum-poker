package websocket

import (
	"encoding/json"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"github.com/scrum-poker/backend/session"
	"log"
	"sync"
	"time"
)

var GlobalHub *Hub

func Init() {
	GlobalHub = newHub()
	go GlobalHub.run()

	session.InitSessionManager(func(roomId string, room *models.Room) {
		GlobalHub.BroadcastRoomUpdate(roomId, room)
	})
}

type RoomData struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

type Hub struct {
	rooms      map[string]*RoomData
	roomsMu    sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*RoomData),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			go h.handleRegistration(client)
		case client := <-h.unregister:
			go h.handleUnregistration(client)
		}
	}
}

func (h *Hub) getOrCreateRoom(roomId string) *RoomData {
	h.roomsMu.RLock()
	roomData, exists := h.rooms[roomId]
	h.roomsMu.RUnlock()

	if !exists {
		h.roomsMu.Lock()
		roomData, exists = h.rooms[roomId]
		if !exists {
			roomData = &RoomData{
				clients: make(map[*Client]bool),
			}
			h.rooms[roomId] = roomData
		}
		h.roomsMu.Unlock()
	}

	return roomData
}

func (h *Hub) handleRegistration(c *Client) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in registration: %v", r)
		}
	}()

	log.Printf("Hub received register request for client %s", c.userId)

	roomData := h.getOrCreateRoom(c.roomId)

	lockAcquired := make(chan bool, 1)
	go func() {
		roomData.mu.Lock()
		lockAcquired <- true
	}()

	select {
	case <-lockAcquired:
		defer roomData.mu.Unlock()

		log.Printf("Processing registration for client %s", c.userId)

		for existingClient := range roomData.clients {
			if existingClient.userId == c.userId {
				log.Printf("Found existing client with userId %s, replacing", c.userId)

				delete(roomData.clients, existingClient)
				close(existingClient.send)

				if existingClient.conn != nil {
					existingClient.conn.Close()
				}
			}
		}

		roomData.clients[c] = true
		log.Printf("Client registered in room %s. Total clients: %d",
			c.roomId, len(roomData.clients))

		select {
		case c.registrationComplete <- true:
		default:
		}

	case <-time.After(5 * time.Second):
		log.Printf("Timeout waiting for lock in registration for client %s", c.userId)
		select {
		case c.registrationComplete <- false:
		default:
		}
	}
}

func (h *Hub) handleUnregistration(c *Client) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in unregistration: %v", r)
		}
	}()

	log.Printf("Hub received unregister request for client %s", c.userId)

	h.roomsMu.RLock()
	roomData, roomExists := h.rooms[c.roomId]
	h.roomsMu.RUnlock()

	if !roomExists {
		log.Printf("Room %s not found for client %s", c.roomId, c.userId)
		return
	}

	lockAcquired := make(chan bool, 1)
	go func() {
		roomData.mu.Lock()
		lockAcquired <- true
	}()

	select {
	case <-lockAcquired:
		defer roomData.mu.Unlock()

		log.Printf("Processing unregistration for client %s", c.userId)

		if _, ok := roomData.clients[c]; ok {
			delete(roomData.clients, c)

			select {
			case c.send <- []byte(`{"type":"disconnect"}`):
			default:
			}

			close(c.send)

			shouldDeleteRoom := len(roomData.clients) == 0
			roomData.mu.Unlock()

			if shouldDeleteRoom {
				h.roomsMu.Lock()
				delete(h.rooms, c.roomId)
				h.roomsMu.Unlock()
				log.Printf("Room %s deleted (no clients)", c.roomId)
			} else {
				log.Printf("Client unregistered from room %s. Total clients: %d",
					c.roomId, len(roomData.clients))
			}

			roomData.mu.Lock()

			go h.handleDisconnection(c)
		}

		log.Printf("Unregistration complete for client %s", c.userId)

	case <-time.After(5 * time.Second):
		log.Printf("Timeout waiting for lock in unregistration for client %s", c.userId)
		if c.conn != nil {
			c.conn.Close()
		}
	}
}

func (h *Hub) handleDisconnection(client *Client) {
	roomId := client.roomId
	userId := client.userId

	go func() {
		room, err := db.GetRoom(roomId)
		if err != nil {
			log.Printf("Error getting room %s: %v", roomId, err)
			return
		}

		if _, ok := room.Participants[userId]; ok {
			if err := db.UpdateUserOnlineStatus(userId, false); err != nil {
				log.Printf("Error updating user online status: %v", err)
				return
			}

			existingSession, err := db.GetSessionByUserID(userId)
			if err == nil && existingSession != nil {
				existingSession.Refresh(session.SessionTTL)
				if err := db.UpdateSession(existingSession); err != nil {
					log.Printf("Error updating session: %v", err)
				}
			} else {
				_, err = session.GlobalManager.CreateSession(userId, roomId)
				if err != nil {
					log.Printf("Error creating session: %v", err)
				}
			}

			updatedRoom, err := db.GetRoom(roomId)
			if err == nil {
				go h.BroadcastRoomUpdate(roomId, updatedRoom)
			}
		}
	}()
}

func (h *Hub) Broadcast(roomId string, message []byte) {
	h.roomsMu.RLock()
	roomData, roomExists := h.rooms[roomId]
	h.roomsMu.RUnlock()

	if !roomExists {
		return
	}

	roomData.mu.Lock()
	defer roomData.mu.Unlock()

	for client := range roomData.clients {
		select {
		case client.send <- message:
		default:
			select {
			case h.unregister <- client:
			default:
				log.Printf("Failed to queue unregister for client %s - channel full", client.userId)
			}
			delete(roomData.clients, client)
		}
	}
}

func (h *Hub) BroadcastRoomUpdate(roomId string, room *models.Room) {
	msg := Message{
		Type:    "room_update",
		RoomId:  roomId,
		Payload: room.ToJSON(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling room update message: %v", err)
		return
	}

	h.Broadcast(roomId, msgBytes)
}
