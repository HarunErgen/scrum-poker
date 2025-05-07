package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string      `json:"type"`
	RoomID  string      `json:"room_id"`
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	hub                  *Hub
	conn                 *websocket.Conn
	send                 chan []byte
	roomID               string
	userID               string
	registrationComplete chan bool
	lastPingTime         time.Time
	mu                   sync.Mutex
}

func newClient(hub *Hub, conn *websocket.Conn, roomID, userID string) *Client {
	return &Client{
		hub:                  hub,
		conn:                 conn,
		send:                 make(chan []byte, 256),
		roomID:               roomID,
		userID:               userID,
		registrationComplete: make(chan bool, 1),
		lastPingTime:         time.Now(),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastPingTime = time.Now()
		c.mu.Unlock()
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		msg.RoomID = c.roomID
		msg.UserID = c.userID

		message, err = json.Marshal(msg)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			continue
		}

		log.Printf("Received message from client %s: %s", c.userID, string(message))

		c.hub.Broadcast(c.roomID, message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	log.Printf("Starting write pump for client %s", c.userID)

	for {
		select {
		case message, ok := <-c.send:
			log.Printf("Writing message to client %s", c.userID)

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("error closing writer: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error pinging client: %v", err)
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, roomID, userID string) {
	log.Printf("Connecting user ID: %s to room %s", userID, roomID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(hub, conn, roomID, userID)

	select {
	case hub.register <- client:
	case <-time.After(2 * time.Second):
		log.Printf("Timeout sending registration request for client %s", userID)
		conn.Close()
		return
	}

	select {
	case success := <-client.registrationComplete:
		if success {
			log.Printf("User ID %s connected to room %s", userID, roomID)

			go client.writePump()
			go client.readPump()
		} else {
			log.Printf("Registration failed for user %s", userID)
			conn.Close()
		}
	case <-time.After(10 * time.Second):
		log.Printf("Registration confirmation timeout for user %s", userID)
		conn.Close()
		return
	}
}
