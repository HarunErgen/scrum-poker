# Scrum Poker Web App

**🔗 Live Demo:** [https://scrumpoker.harunergen.com/](https://scrumpoker.harunergen.com/)

Scrum Poker is a real-time estimation tool built for Agile teams to collaboratively plan sprints. Participants can create or join rooms, submit task estimates using the Fibonacci sequence, and reveal votes in sync—facilitating transparent, fast-paced planning meetings.

---

## 🧱 Project Architecture

This project consists of two main components:

* **Frontend:** A responsive React.js app (mobile + desktop)
* **Backend:** A Go service with WebSocket support and PostgreSQL for persistence

---

## 🚀 Features

* Create or join real-time estimation rooms
* Support for multiple concurrent rooms
* Scrum Master role with voting control
* Fibonacci-based vote selection
* Vote reveal/reset functionality
* Automatic user presence updates (online/offline)
* Resilient to short-term disconnections

---

## 🕸️ Real-Time Communication

Scrum Poker uses WebSockets for persistent, low-latency messaging. Here's how the system handles real-time collaboration:

### 📌 Persistent Connections

Each client connects once and maintains a live WebSocket session:

![1](https://github.com/user-attachments/assets/f6171397-020a-4810-a70d-3a8226fe1d35)

### 🔁 Message Broadcasting

Whenever a user performs an action (e.g., submit vote, rename), a structured WebSocket message is sent to the backend. It processes the request and rebroadcasts an updated room state to all clients:

![2 (1)](https://github.com/user-attachments/assets/7cf88456-b51e-4454-a612-1eec1665a65d)

Example message structure:

```go
Message {
  Action: models.ActionType,
  Payload: map[string]interface{}
}
```

---

### 🧠 Backend Message Handling (Go)

Each WebSocket message is processed by the backend's `ProcessMessage` function. Actions like `submit`, `reveal`, `rename`, or `leave` are handled explicitly:

```go
func ProcessMessage(broadcastFunc, roomId, msg)
```

Example: Submitting a vote

```go
Payload: {
  "userId": "123",
  "vote": "5"
}
```

The backend updates the database and broadcasts:

```go
Action: "submit",
Payload: { "userId": "123", "vote": "5" }
```

---

### 🎯 Frontend Message Handling (React)

Messages are processed on the UI with a centralized handler:

```js
processWebSocketMessage(message, roomData, setRoomData, setSelectedVote)
```

This dynamically updates room state based on action types like:

* `join` → Add user to participants
* `submit` → Register vote
* `reveal` → Reveal all votes
* `reset` → Clear votes and reset selection
* `rename`, `leave`, `offline`, `online`, `transfer` → Reflect changes instantly

UI feedback (e.g., toasts) and state re-renders are triggered accordingly.

---

## ☁️ Deployment (EC2)

The Scrum Poker app is deployed on a single **Amazon EC2 instance**.

* All services (frontend, backend, and database) run as Docker containers managed by Docker Compose.
* The **React frontend** is served via Nginx and exposed over HTTPS.
* The **Go backend** handles API and WebSocket traffic.
* The **PostgreSQL** database is used for persistence.

### 🔁 Request/Response Flow

```plaintext
       User Browser
         (Client)
            │
         HTTPS / WSS
            │
            ▼
    ┌────────────────────┐
    │    EC2 Instance    │
    │  (Docker Compose)  │
    └────────────────────┘
        │         │
        ▼         ▼
  ┌────────┐   ┌────────┐
  │Frontend│   │Backend │◄───▶ PostgreSQL
  │ React  │   │   Go   │
  └────────┘   └────────┘
        ▲         ▲
        └─────────┘
       Internal Network
```

* The browser sends HTTP(S) and WebSocket requests to the EC2 public IP or domain.
* Nginx (inside the frontend container) serves the React app and proxies API/WebSocket requests to the backend.
* The backend processes logic and communicates with the PostgreSQL container.
* Responses and real-time updates are sent back to clients via WebSocket or HTTP.

---

## 🧑‍💻 Local Development

### Frontend

* **Requirements:** Node.js 16+, npm or yarn

```bash
cd frontend
npm install
npm start
```

**Build for production:**

```bash
npm run build
```

**Env Variables:**

* `REACT_APP_API_URL` – Backend endpoint

### Backend

* **Requirements:** Go 1.23+, PostgreSQL 14+

```bash
cd backend
go mod download
go run main.go
```

**Env Variables:**

* `BACKEND_PORT`
* `FRONTEND_PORT`
* `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
* `REACT_APP_API_URL`
* `ALLOWED_ORIGINS` (for CORS)

---

## 📡 API & WebSocket

### REST API

| Action           | Method | Endpoint               |
| ---------------- | ------ | ---------------------- |
| Create Room      | POST   | `/rooms`               |
| Join Room        | POST   | `/rooms/{roomId}/join` |
| Get Room Details | GET    | `/rooms/{roomId}`      |

### WebSocket

| Endpoint                       | Description                    |
| ------------------------------ | ------------------------------ |
| `/ws/{roomId}?userId={userId}` | Establish WebSocket connection |

**Example Message Payload:**

```json
{
  "action": "submit",
  "payload": {
    "userId": "abc123",
    "vote": "5"
  }
}
```

All messages follow this pattern and are rebroadcast to clients for UI synchronization.

---

## 🔐 Session Management

Session management is cookie-based and supports short disconnections:

| Action         | Method | Endpoint                         |
| -------------- | ------ | -------------------------------- |
| Create Session | POST   | `/sessions/{userId}/{roomId}`    |
| Validate       | GET    | `/sessions?roomId={roomId}`      |
| Delete         | DELETE | `/sessions` (via session cookie) |

---

## ⚙️ Tech Stack

### Backend

* Language: **Go**
* Database: **PostgreSQL**
* Real-time: **WebSockets**
* API: **REST (JSON)**

### Frontend

* Framework: **React.js**
* State: **React Hooks**
* Routing: **React Router**
* HTTP: **Axios**
* Served via **Nginx**

---

## 🐳 Dockerized Environment

All services are containerized and orchestrated via Docker Compose:

* `frontend` – React app served through Nginx
* `backend` – Go-based WebSocket server
* `postgres` – Persistent database

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
