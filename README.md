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
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
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

- Go 1.23+
- PostgreSQL (or use the Docker Compose setup)

#### Installation

```bash
# Navigate to backend directory
cd backend

# Download dependencies
go mod download

# Run the application
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

## API Documentation

### Room Management

#### Create a Room

- **URL**: `/api/rooms`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Sprint Planning",
    "user_name": "John Doe"
  }
  ```
- **Response**:
  ```json
  {
    "id": "room-id",
    "name": "Sprint Planning",
    "created_at": "2023-01-01T12:00:00Z",
    "scrum_master": "user-id",
    "participants": {
      "user-id": {
        "id": "user-id",
        "name": "John Doe",
        "created_at": "2023-01-01T12:00:00Z"
      }
    },
    "votes": {},
    "votes_revealed": false
  }
  ```

#### Get Room Details

- **URL**: `/api/rooms/{roomID}`
- **Method**: `GET`
- **Response**: Room details (same format as above)

#### Join a Room

- **URL**: `/api/rooms/{roomID}/join`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_name": "Jane Smith"
  }
  ```
- **Response**:
  ```json
  {
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "created_at": "2023-01-01T12:00:00Z",
      "scrum_master": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "created_at": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votes_revealed": false
    },
    "user": {
      "id": "user-id",
      "name": "Jane Smith",
      "created_at": "2023-01-01T12:05:00Z"
    }
  }
  ```

### Voting

#### Submit a Vote

- **URL**: `/api/rooms/{roomID}/vote`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": "user-id",
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
      "created_at": "2023-01-01T12:00:00Z",
      "scrum_master": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "created_at": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votes_revealed": false
    }
  }
  ```

#### Reveal Votes

- **URL**: `/api/rooms/{roomID}/reveal`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": "scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "created_at": "2023-01-01T12:00:00Z",
      "scrum_master": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "created_at": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {
        "user-id": "5"
      },
      "votes_revealed": true
    }
  }
  ```

#### Reset Votes

- **URL**: `/api/rooms/{roomID}/reset`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": "scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "created_at": "2023-01-01T12:00:00Z",
      "scrum_master": "user-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "created_at": "2023-01-01T12:00:00Z"
        }
      },
      "votes": {},
      "votes_revealed": false
    }
  }
  ```

#### Transfer Scrum Master Role

- **URL**: `/api/rooms/{roomID}/transfer`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": "current-scrum-master-id",
    "new_scrum_master_id": "new-scrum-master-id"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "room": {
      "id": "room-id",
      "name": "Sprint Planning",
      "created_at": "2023-01-01T12:00:00Z",
      "scrum_master": "new-scrum-master-id",
      "participants": {
        "user-id": {
          "id": "user-id",
          "name": "John Doe",
          "created_at": "2023-01-01T12:00:00Z"
        },
        "new-scrum-master-id": {
          "id": "new-scrum-master-id",
          "name": "Jane Smith",
          "created_at": "2023-01-01T12:05:00Z"
        }
      },
      "votes": {},
      "votes_revealed": false
    }
  }
  ```

### WebSocket

- **URL**: `/ws/{roomID}?user_id={userID}`
- **Protocol**: WebSocket
- **Messages**:
  - Room updates are sent as JSON with the following format:
    ```json
    {
      "type": "room_update",
      "room_id": "room-id",
      "user_id": "user-id",
      "payload": {
        "id": "room-id",
        "name": "Sprint Planning",
        "created_at": "2023-01-01T12:00:00Z",
        "scrum_master": "user-id",
        "participants": {
          "user-id": {
            "id": "user-id",
            "name": "John Doe",
            "created_at": "2023-01-01T12:00:00Z"
          }
        },
        "votes": {},
        "votes_revealed": false
      }
    }
    ```

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