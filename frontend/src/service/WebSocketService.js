import Message from "../models/Message";

class WebSocketService {
  constructor(roomId, userId, onMessage, onStatusChange) {
    this.roomId = roomId;
    this.userId = userId;
    this.onMessage = onMessage;
    this.onStatusChange = onStatusChange;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
    this.reconnectTimeout = null;
    this.heartbeatInterval = null;
    this.lastPongTime = null;
    this.messageQueue = [];
    this.isConnected = false;
  }

  connect() {
    if (this.isConnected) {
      console.log('Already connected, not connecting again');
      return;
    }

    if (this.ws) {
      console.log('Existing WebSocket found, disconnecting first');
      this.disconnect();
    }

    if (!this.userId) {
      console.error('Cannot connect: userId is null or undefined');
      this.onStatusChange('Error: Missing user ID');

      if (!this.reconnectTimeout) {
        console.log('Scheduling reconnect due to missing userId');
        this.reconnectAttempts = 0;
        setTimeout(() => this.reconnect(), 2000);
      } else {
        console.log('Reconnect already scheduled, not scheduling another');
      }
      return;
    }

    console.log(`Connecting WebSocket with roomId=${this.roomId} and userId=${this.userId}`);
    this.ws = new WebSocket(`${process.env.REACT_APP_API_URL}/ws/${this.roomId}?userId=${this.userId}`);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.isConnected = true;
      this.reconnectAttempts = 0;
      this.onStatusChange('Connected');

      this.startHeartbeat();
      this.flushMessageQueue();
    };

    this.ws.onmessage = async (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.action === 'pong') {
          if (data.payload && data.payload.userId === this.userId) {
            console.log('Received pong for our ping');
            this.lastPongTime = Date.now();
          }
          return;
        }

        this.lastPongTime = Date.now();

        const message = new Message(data.action, data.payload);
        await this.onMessage(message);
      } catch (error) {
        console.error('Error processing message:', error);
      }
    };

    this.ws.onclose = (event) => {
      console.log(`WebSocket disconnected with code: ${event.code}, reason: ${event.reason || 'No reason provided'}`);
      this.isConnected = false;
      this.onStatusChange('Disconnected');
      this.stopHeartbeat();

      if (event.code === 1000) {
        console.log('Clean WebSocket close, not reconnecting');
      } else if (this.reconnectTimeout) {
        console.log('Reconnect already scheduled, not scheduling another');
      } else {
        console.log('Unclean WebSocket close, attempting to reconnect');
        this.reconnect();
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.onStatusChange('Error');
    };
  }

  startHeartbeat() {
    this.stopHeartbeat();

    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        console.log('Sending ping');

        try {
          const pingMessage = new Message('ping', { userId: this.userId });
          this.ws.send(JSON.stringify(pingMessage));

          if (this.lastPongTime && Date.now() - this.lastPongTime > 30000) {
            console.log('No pong received in 30 seconds, reconnecting...');
            this.disconnect();
            this.reconnect();
          }
        } catch (error) {
          console.error('Error sending ping:', error);
        }
      }
    }, 15000);

    this.lastPongTime = Date.now();
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  reconnect() {
    if (this.reconnectTimeout) {
      console.log('Clearing existing reconnect timeout');
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }

    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnect attempts reached');
      this.onStatusChange('Disconnected - Max retries reached');
      return;
    }

    if (this.isConnected) {
      console.log('Already connected, skipping reconnect');
      return;
    }

    this.reconnectAttempts++;

    const baseDelay = Math.min(Math.pow(2, this.reconnectAttempts - 1) * 1000, 30000);
    const jitter = baseDelay * 0.2 * (Math.random() - 0.5);
    const delay = baseDelay + jitter;

    console.log(`Reconnecting in ${Math.round(delay/1000)}s (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
    this.onStatusChange(`Reconnecting in ${Math.round(delay/1000)}s...`);

    this.reconnectTimeout = setTimeout(() => {
      if (!this.isConnected) {
        console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        this.reconnectTimeout = null;
        this.connect();
      } else {
        console.log('Connection already established, canceling reconnect attempt');
        this.reconnectTimeout = null;
      }
    }, delay);
  }

  flushMessageQueue() {
    if (this.messageQueue.length > 0) {
      console.log(`Flushing ${this.messageQueue.length} queued messages`);

      const queueCopy = [...this.messageQueue];
      this.messageQueue = [];

      queueCopy.forEach(msg => {
        this.sendMessage(msg.action, msg.payload);
      });
    }
  }

  sendMessage(action, payload) {
    const isPing = action === 'ping';

    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');

      if (!isPing && action !== 'pong') {
        console.log(`Queuing message: ${action}`);
        this.messageQueue.push({ action, payload });
      }

      return false;
    }

    try {
      const message = new Message(action, payload);
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error('Error sending WebSocket message:', error);

      if (!isPing && action !== 'pong') {
        console.log(`Queuing message after error: ${action}`);
        this.messageQueue.push({ action, payload });
      }

      return false;
    }
  }

  disconnect() {
    console.log('Disconnecting WebSocket');

    this.stopHeartbeat();

    if (this.reconnectTimeout) {
      console.log('Clearing reconnect timeout during disconnect');
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }

    if (this.ws) {
      try {
        if (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING) {
          console.log('Closing WebSocket connection');
          this.ws.close(1000, 'Normal closure');
        } else {
          console.log(`WebSocket already in ${this.ws.readyState === WebSocket.CLOSING ? 'CLOSING' : 'CLOSED'} state`);
        }
      } catch (error) {
        console.error('Error closing WebSocket:', error);
      }

      this.ws.onopen = null;
      this.ws.onmessage = null;
      this.ws.onclose = null;
      this.ws.onerror = null;

      this.ws = null;
    }

    this.isConnected = false;
    console.log('WebSocket disconnected');
  }

  updateUserId(newUserId) {
    console.log(`updateUserId called with newUserId=${newUserId}, current userId=${this.userId}, isConnected=${this.isConnected}, hasReconnectTimeout=${!!this.reconnectTimeout}`);

    if (!newUserId) {
      console.warn('Attempted to update userId with null or undefined value');
      return;
    }

    if (this.userId !== newUserId) {
      console.log(`Updating userId from ${this.userId} to ${newUserId}`);
      this.userId = newUserId;

      if (this.reconnectTimeout) {
        console.log('Canceling existing reconnect timeout to start fresh with new userId');
        clearTimeout(this.reconnectTimeout);
        this.reconnectTimeout = null;
        this.reconnectAttempts = 0;
      }

      if (!this.isConnected) {
        console.log('WebSocket is disconnected. Connecting now with new userId.');
        this.connect();
      } else {
        console.log('WebSocket is already connected. Will use new userId on next reconnect.');
      }
    } else {
      console.log(`UserId unchanged (${this.userId})`);
    }
  }
}

export default WebSocketService;
