class Room {
  constructor(data) {
    this.id = data.id || '';
    this.name = data.name || '';
    this.scrumMaster = data.scrum_master || '';
    this.participants = data.participants || {};
    this.votes = data.votes || {};
    this.votesRevealed = data.votes_revealed || false;
  }

  isScrumMaster(userId) {
    return this.scrumMaster === userId;
  }

  hasVoted(userId) {
    return !!this.votes[userId];
  }

  getUserVote(userId) {
    return this.votes[userId];
  }

  getParticipantsArray() {
    return Object.values(this.participants);
  }

  static fromApiResponse(data) {
    return new Room(data);
  }
}

export default Room;