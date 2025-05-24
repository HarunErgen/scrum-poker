import Message from "../models/Message";

class WebSocketService {
  constructor(roomId, userId, onMessage, onStatusChange) {
    this.roomId = roomId;
    this.userId = userId;
    this.onMessage = onMessage;
    this.onStatusChange = onStatusChange;
    this.ws = null;
  }

  connect() {
    if (this.ws) {
      this.disconnect();
    }

    this.ws = new WebSocket(`${process.env.REACT_APP_API_URL}/ws/${this.roomId}?userId=${this.userId}`);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.onStatusChange('Connected');
    };

    this.ws.onmessage = async (event) => {
      const data = JSON.parse(event.data);
      const message = new Message(data.action, data.payload);
      await this.onMessage(message);
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.onStatusChange('Disconnected');
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.onStatusChange('Error');
    };
  }

  sendMessage(action, payload) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return false;
    }

    try {
      const message = new Message(action, payload);
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error('Error sending WebSocket message:', error);
      return false;
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export default WebSocketService;
