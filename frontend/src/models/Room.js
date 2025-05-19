class Room {
  constructor(id, name, scrumMaster, participants, votes, votesRevealed) {
    this.id = id || '';
    this.name = name || '';
    this.scrumMaster = scrumMaster || '';
    this.participants = participants || {};
    this.votes = votes || {};
    this.votesRevealed = votesRevealed || false;
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
    return new Room(data.id, data.name, data.scrumMaster, data.participants, data.votes, data.votesRevealed);
  }
}

export default Room;