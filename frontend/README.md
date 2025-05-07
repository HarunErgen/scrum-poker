# Scrum Poker Frontend

This is the frontend implementation for the Scrum Poker Web App. It's built with React.js and communicates with the backend API via HTTP and WebSockets for real-time updates.

## Features

- Create or join estimation rooms
- Vote using Fibonacci sequence cards
- Real-time updates via WebSockets
- Scrum Master controls for revealing and resetting votes
- Responsive design for desktop and mobile

## Development

### Prerequisites

- Node.js 16+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm start
```

### Building for Production

```bash
npm run build
```

## Docker

The application is containerized using Docker. The Dockerfile uses a multi-stage build process:

1. Build stage: Uses Node.js to build the React application
2. Production stage: Uses nginx to serve the static files

The nginx configuration proxies API and WebSocket requests to the backend service.

## Environment Variables

- `REACT_APP_API_URL` - URL of the backend API (default: http://localhost:8080)