class SessionStorage {
  static KEYS = {
    USER_ID: 'userId',
    USER_NAME: 'userName',
    ROOM_ID: 'roomId'
  };

  static saveSession(userId, userName, roomId) {
    localStorage.setItem(this.KEYS.USER_ID, userId);
    localStorage.setItem(this.KEYS.USER_NAME, userName);
    localStorage.setItem(this.KEYS.ROOM_ID, roomId);
  }

  static getSession() {
    return {
      userId: localStorage.getItem(this.KEYS.USER_ID),
      userName: localStorage.getItem(this.KEYS.USER_NAME),
      roomId: localStorage.getItem(this.KEYS.ROOM_ID)
    };
  }

  static hasSession() {
    const { userId, userName, roomId } = this.getSession();
    return !!(userId && userName && roomId);
  }

  static clearSession() {
    localStorage.removeItem(this.KEYS.USER_ID);
    localStorage.removeItem(this.KEYS.USER_NAME);
    localStorage.removeItem(this.KEYS.ROOM_ID);
  }
}

export default SessionStorage;
