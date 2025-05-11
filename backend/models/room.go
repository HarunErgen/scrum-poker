package models

import (
	"sync"
	"time"
)

type Room struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	CreatedAt     time.Time         `json:"created_at"`
	ScrumMaster   string            `json:"scrum_master"`
	Participants  map[string]*User  `json:"participants"`
	Votes         map[string]string `json:"votes"`
	VotesRevealed bool              `json:"votes_revealed"`
	Mu            sync.Mutex        `json:"-"`
}

func NewRoom(id, name, scrumMasterID string) *Room {
	return &Room{
		ID:            id,
		Name:          name,
		CreatedAt:     time.Now(),
		ScrumMaster:   scrumMasterID,
		Participants:  make(map[string]*User),
		Votes:         make(map[string]string),
		VotesRevealed: false,
	}
}

func (r *Room) AddParticipant(user *User) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Participants[user.ID] = user
}

func (r *Room) RemoveParticipant(userID string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Participants, userID)
	delete(r.Votes, userID)
}

func (r *Room) AddVote(userID, vote string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Votes[userID] = vote
}

func (r *Room) RevealVotes() {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.VotesRevealed = true
}

func (r *Room) ResetVotes() {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Votes = make(map[string]string)
	r.VotesRevealed = false
}

func (r *Room) TransferScrumMaster(newScrumMasterID string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ScrumMaster = newScrumMasterID
}

func (r *Room) AssignRandomScrumMaster(participants map[string]*User) {
	var candidates []string
	for userID := range participants {
		candidates = append(candidates, userID)
	}

	if len(candidates) == 0 {
		return
	}

	randomIndex := time.Now().UnixNano() % int64(len(candidates))
	newScrumMasterId := candidates[randomIndex]
	r.TransferScrumMaster(newScrumMasterId)
}

func (r *Room) ToJSON() map[string]interface{} {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	participants := make(map[string]interface{})
	for id, user := range r.Participants {
		participants[id] = user.ToJSON()
	}

	votes := make(map[string]string)
	if r.VotesRevealed {
		votes = r.Votes
	} else {
		for id := range r.Votes {
			votes[id] = "voted"
		}
	}

	return map[string]interface{}{
		"id":             r.ID,
		"name":           r.Name,
		"created_at":     r.CreatedAt,
		"scrum_master":   r.ScrumMaster,
		"participants":   participants,
		"votes":          votes,
		"votes_revealed": r.VotesRevealed,
	}
}

func (r *Room) RemoveVote(userID string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Votes, userID)
}
