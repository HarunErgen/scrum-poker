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

    this.ws = new WebSocket(`${process.env.REACT_APP_API_URL}/ws/${this.roomId}?user_id=${this.userId}`);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.onStatusChange('Connected');
    };

    this.ws.onmessage = async (event) => {
      const data = JSON.parse(event.data);
      await this.onMessage(data);
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

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export default WebSocketService;
