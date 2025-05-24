# Scrum Poker Web App

A real-time Scrum Poker application for Agile teams to efficiently estimate tasks. This application allows teams to create estimation rooms, join existing rooms, and vote on tasks using the Fibonacci sequence.

## Project Structure

The project is organized into two main components:

- **Frontend**: A React.js application with a responsive design for desktop and mobile
- **Backend**: A Go application with WebSocket support for real-time communication and PostgreSQL for data persistence

## Features

- Create or join estimation rooms instantly
- Support for multiple concurrent rooms
- Scrum Master role management
- Real-time voting using Fibonacci sequence
- Vote revelation controlled by the Scrum Master
- Responsive design for desktop and mobile devices
- Session management for short-time disconnections.

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running the Application

1. Clone the repository
2. Run the application using Docker Compose:

```bash
docker-compose up -d
```

3. Access the application:
   - Frontend: http://localhost
   - Backend API: http://localhost:8080/api/health

## Development

### Frontend Development

#### Prerequisites

- Node.js 16+
- npm or yarn

#### Installation

```bash
cd frontend
npm install
npm start
```

#### Building for Production

```bash
npm run build
```

#### Environment Variables

- `REACT_APP_API_URL` - URL of the backend API (default: http://localhost:8080)

### Backend Development

#### Prerequisites

- Go 1.23.0
- PostgreSQL 14 (or use the Docker Compose setup)

#### Installation

```bash
cd backend
go mod download
go run main.go
```

#### Environment Variables

- `BACKEND_PORT` - Server port (default: 8080)
- `FRONTEND_PORT` - Port for the frontend application (default: 80)
- `REACT_APP_API_URL` - API url (default: http://localhost:8080)
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - PostgreSQL user
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name
- `DB_SSLMODE` - PostgreSQL SSL mode
- `ALLOWED_ORIGINS` - Allowed CORS origins.

## API Documentation

### User Management
#### Update User

- **URL**: `/api/users`
- **Method**: `PUT`
- **Request Body**:
  ```json
  {
    "id": "user-id",
    "name": "John Doe",
    "isOnline": true
  }
  ```
- **Response**:
  ```json
  {
    "id": "user-id",
    "name": "John Doe",
    "isOnline": true,
    "createdAt": "2023-01-01T12:00:00Z"
  }
  ```

### Room Management

#### Create a Room

- **URL**: `/api/rooms`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Sprint Planning",
    "userName": "John Doe"
  }
  ```
- **Response**:
  ```json
  {
    "id": "room-id",
    "name": "Sprint Planning",
    "createdAt": "2023-01-01T12:00:00Z",
    "scrumMaster": "user-id",
    "participants": {
      "user-id": {
        "id": "user-id",
        "name": "John Doe",
        "isOnline": true,
        "createdAt": "2023-01-01T12:00:00Z"
      }
    },
    "votes": {},
    "votesRevealed": false
  }
  ```

#### Get Room Details

- **URL**: `/api/rooms/{roomId}`
- **Method**: `GET`
- **Response**: Room details (same format as above)

#### Join a Room

- **URL**: `/api/rooms/{roomId}/join`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "userName": "Jane Smith"
  }
  ```
- **Response**:
  ```json
  {
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "createdAt": "2023-01-01T12:00:00Z",
      "scrumMaster": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "isOnline": true,
          "createdAt": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votesRevealed": false
    },
    "user": {
      "id": "user-id",
      "name": "Jane Smith",
      "isOnline": true,
      "createdAt": "2023-01-01T12:05:00Z"
    }
  }
  ```

### Voting

#### Submit a Vote

- **URL**: `/api/rooms/{roomId}/vote`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "userId": "user-id",
    "vote": "5"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "createdAt": "2023-01-01T12:00:00Z",
      "scrumMaster": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "isOnline": true,
          "createdAt": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votesRevealed": false
    }
  }
  ```

#### Reveal Votes

- **URL**: `/api/rooms/{roomId}/reveal`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "userId": "scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "createdAt": "2023-01-01T12:00:00Z",
      "scrumMaster": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "isOnline": true,
          "createdAt": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {
        "user-id": "5"
      },
      "votesRevealed": true
    }
  }
  ```

#### Reset Votes

- **URL**: `/api/rooms/{roomId}/reset`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "userId": "scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "createdAt": "2023-01-01T12:00:00Z",
      "scrumMaster": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "isOnline": true,
          "createdAt": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votesRevealed": false
    }
  }
  ```

#### Transfer Scrum Master Role

- **URL**: `/api/rooms/{roomId}/transfer`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "userId": "current-scrum-master-id",
    "newScrumMasterId": "new-scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "createdAt": "2023-01-01T12:00:00Z",
      "scrumMaster": "new-scrum-master-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "isOnline": true,
          "createdAt": "2023-01-01T12:00:00Z"
        },
        "new-scrum-master-id": {
          "id": "new-scrum-master-id",
          "name": "Jane Smith",
          "isOnline": true,
          "createdAt": "2023-01-01T12:05:00Z"
        }
      },
      "votes": {},
      "votesRevealed": false
    }
  }
  ```

### WebSocket

- **URL**: `/ws/{roomId}?userId={userId}`
- **Protocol**: WebSocket
- **Messages**:
  - Room updates are sent as JSON with the following format:
    ```json
    {
      "type": "room_update",
      "roomId": "room-id",
      "userId": "user-id",
      "payload": {
        "id": "room-id",
        "name": "Sprint Planning",
        "createdAt": "2023-01-01T12:00:00Z",
        "scrumMaster": "user-id",
        "participants": {
          "user-id": {
            "id": "user-id",
            "name": "John Doe",
            "isOnline": true,
            "createdAt": "2023-01-01T12:00:00Z"
          }
        },
        "votes": {},
        "votesRevealed": false
      }
    }
    ```

### Session Management

#### Create Session

- **URL**: `/api/sessions/{userId}/{roomId}`
- **Method**: `POST`
- **Request Body**: `None`

- **Response**:
  - Status: `201 Created`
  - Set-Cookie: `sessionId=<session-id>; HttpOnly; Secure; SameSite=Lax`
  - Body:
      ```json
    {
        "sessionId": "session-id"
    }
    ```

- **Error Responses**:
  - `400 Bad Request`: If userId or roomId is missing.
  - `403 Forbidden`: If the user is not a participant in the room.
  - `404 Not Found`: If the user or room does not exist.
  - `500 Internal Server Error`: If session creation or user status update fails.

---

#### Get Session
- **URL**: `/api/sessions`
- **Method**: `GET`
- **Query Parameters**:
  - `roomId` (string) – The ID of the room to validate against the session.

- **Request Headers**:
  - `Cookie`: sessionId=<session-id>

- **Response**:
  - Status: 200 OK
  - Body:
  ```json
    {
      "session": {
        "id": "session-id",
        "userId": "user-id",
        "roomId": "room-id",
        "createdAt": "2023-01-01T12:00:00Z",
        "expiresAt": "2023-01-01T13:00:00Z"
      },
      "user": {
        "id": "user-id",
        "name": "John Doe",
        "isOnline": true,
        "createdAt": "2023-01-01T12:00:00Z"
      },
      "room": {
        "id": "room-id",
        "name": "Sprint Planning",
        "createdAt": "2023-01-01T12:00:00Z",
        "scrumMaster": "user-id",
        "participants": {
          "user-id": {
            "id": "user-id",
            "name": "John Doe",
            "isOnline": true,
            "createdAt": "2023-01-01T12:00:00Z"
          }
        },
        "votes": {},
        "votesRevealed": false
      }
    }
  ```

- **Error Responses**:
  - `400 Bad Request`: If the session ID is missing.
  - `401 Unauthorized`: If the session has expired.
  - `403 Forbidden`: If the room ID in query does not match the session.
  - `404 Not Found`: If the session, user, or room is not found.
  - `500 Internal Server Error`: If session refresh or user update fails.

---

#### Delete Session

- **URL**: `/api/sessions`
- **Method**: `DELETE`
- **Request Headers**:
  - Cookie: sessionId=<session-id>
- **Response**:
  - Status: `200 OK`
  - Set-Cookie: `sessionId=; Max-Age=-1; HttpOnly; Secure; SameSite=Lax`
  - Body:
    ```json
    {
      "message": "Session deleted"
    }
    ```

- Error Responses:
  - `400 Bad Request`: If session ID is missing.
  - `404 Not Found`: If the session cookie is not present.
  - `500 Internal Server Error`: If session deletion fails.

## Technical Details

### Backend

- **Language**: Go
- **Database**: PostgreSQL
- **Real-time Communication**: WebSockets
- **API**: RESTful with JSON

### Frontend

- **Framework**: React.js
- **HTTP Client**: Axios
- **Routing**: React Router
- **Styling**: CSS
- **Server**: Nginx (in production)

## Docker Configuration

The application is fully dockerized with the following services:

- **backend**: Go application
- **postgres**: PostgreSQL database
- **frontend**: React application served via Nginx

## License

MIT
