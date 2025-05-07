package websocket

import (
	"encoding/json"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/models"
	"log"
	"sync"
	"time"
)

var CommonHub *Hub

func Init() {
	CommonHub = newHub()
	go CommonHub.run()
	go CommonHub.startHealthCheck()
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

func (h *Hub) getOrCreateRoom(roomID string) *RoomData {
	h.roomsMu.RLock()
	roomData, exists := h.rooms[roomID]
	h.roomsMu.RUnlock()

	if !exists {
		h.roomsMu.Lock()
		roomData, exists = h.rooms[roomID]
		if !exists {
			roomData = &RoomData{
				clients: make(map[*Client]bool),
			}
			h.rooms[roomID] = roomData
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

	log.Printf("Hub received register request for client %s", c.userID)

	roomData := h.getOrCreateRoom(c.roomID)

	lockAcquired := make(chan bool, 1)
	go func() {
		roomData.mu.Lock()
		lockAcquired <- true
	}()

	select {
	case <-lockAcquired:
		defer roomData.mu.Unlock()

		log.Printf("Processing registration for client %s", c.userID)

		for existingClient := range roomData.clients {
			if existingClient.userID == c.userID {
				log.Printf("Found existing client with userID %s, replacing", c.userID)

				delete(roomData.clients, existingClient)
				close(existingClient.send)

				if existingClient.conn != nil {
					existingClient.conn.Close()
				}
			}
		}

		roomData.clients[c] = true
		log.Printf("Client registered in room %s. Total clients: %d",
			c.roomID, len(roomData.clients))

		select {
		case c.registrationComplete <- true:
		default:
		}

	case <-time.After(5 * time.Second):
		log.Printf("Timeout waiting for lock in registration for client %s", c.userID)
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

	log.Printf("Hub received unregister request for client %s", c.userID)

	h.roomsMu.RLock()
	roomData, roomExists := h.rooms[c.roomID]
	h.roomsMu.RUnlock()

	if !roomExists {
		log.Printf("Room %s not found for client %s", c.roomID, c.userID)
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

		log.Printf("Processing unregistration for client %s", c.userID)

		if _, ok := roomData.clients[c]; ok {
			delete(roomData.clients, c)

			select {
			case c.send <- []byte(`{"type":"disconnect"}`):
			default:
			}

			close(c.send)

			if len(roomData.clients) == 0 {
				h.roomsMu.Lock()
				delete(h.rooms, c.roomID)
				h.roomsMu.Unlock()
				log.Printf("Room %s deleted (no clients)", c.roomID)
			} else {
				log.Printf("Client unregistered from room %s. Total clients: %d",
					c.roomID, len(roomData.clients))
			}

			go h.handleDisconnection(c)
		}

		log.Printf("Unregistration complete for client %s", c.userID)

	case <-time.After(5 * time.Second):
		log.Printf("Timeout waiting for lock in unregistration for client %s", c.userID)
		if c.conn != nil {
			c.conn.Close()
		}
	}
}

func (h *Hub) handleDisconnection(client *Client) {
	roomID := client.roomID
	userID := client.userID

	go func() {
		room, err := db.GetRoom(roomID)
		if err != nil {
			log.Printf("Error getting room %s: %v", roomID, err)
			return
		}

		if user, ok := room.Participants[userID]; ok {
			user.IsActive = false

			_, err = db.DB.Exec(
				"UPDATE users SET is_active = $1 WHERE id = $2",
				false, userID,
			)
			if err != nil {
				log.Printf("Error updating user status: %v", err)
				return
			}

			go h.BroadcastRoomUpdate(roomID, room)
		}
	}()
}

func (h *Hub) Broadcast(roomID string, message []byte) {
	h.roomsMu.RLock()
	roomData, roomExists := h.rooms[roomID]
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
				log.Printf("Failed to queue unregister for client %s - channel full", client.userID)
			}
			delete(roomData.clients, client)
		}
	}
}

func (h *Hub) BroadcastRoomUpdate(roomID string, room *models.Room) {
	msg := Message{
		Type:    "room_update",
		RoomID:  roomID,
		Payload: room.ToJSON(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling room update message: %v", err)
		return
	}

	h.Broadcast(roomID, msgBytes)
}

func (h *Hub) startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		h.checkConnections()
	}
}

func (h *Hub) checkConnections() {
	now := time.Now()
	timeout := 2 * pongWait

	h.roomsMu.RLock()
	roomsToCheck := make([]string, 0, len(h.rooms))
	for roomID := range h.rooms {
		roomsToCheck = append(roomsToCheck, roomID)
	}
	h.roomsMu.RUnlock()

	for _, roomID := range roomsToCheck {
		h.roomsMu.RLock()
		roomData, roomExists := h.rooms[roomID]
		h.roomsMu.RUnlock()

		if !roomExists {
			continue
		}

		var clientsToRemove []*Client

		roomData.mu.Lock()
		for client := range roomData.clients {
			client.mu.Lock()
			lastPing := client.lastPingTime
			client.mu.Unlock()

			if now.Sub(lastPing) > timeout {
				clientsToRemove = append(clientsToRemove, client)
			}
		}
		roomData.mu.Unlock()

		for _, client := range clientsToRemove {
			log.Printf("Client %s has timed out, removing", client.userID)

			select {
			case h.unregister <- client:
			default:
				if client.conn != nil {
					client.conn.Close()
				}
			}
		}
	}
}
