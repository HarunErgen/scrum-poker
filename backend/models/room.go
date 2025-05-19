package models

import (
	"sync"
	"time"
)

type Room struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	CreatedAt     time.Time         `json:"createdAt"`
	ScrumMaster   string            `json:"scrumMaster"`
	Participants  map[string]*User  `json:"participants"`
	Votes         map[string]string `json:"votes"`
	VotesRevealed bool              `json:"votesRevealed"`
	Mu            sync.Mutex        `json:"-"`
}

func NewRoom(id, name, scrumMasterID string) *Room {
	return &Room{
		Id:            id,
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
	r.Participants[user.Id] = user
}

func (r *Room) RemoveParticipant(userId string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Participants, userId)
	delete(r.Votes, userId)
}

func (r *Room) AddVote(userId, vote string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Votes[userId] = vote
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
	for userId := range participants {
		candidates = append(candidates, userId)
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
		"id":            r.Id,
		"name":          r.Name,
		"createdAt":     r.CreatedAt,
		"scrumMaster":   r.ScrumMaster,
		"participants":  participants,
		"votes":         votes,
		"votesRevealed": r.VotesRevealed,
	}
}

func (r *Room) RemoveVote(userId string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Votes, userId)
}
