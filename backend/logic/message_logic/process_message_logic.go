package message_logic

import (
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/logic/room_logic"
	"github.com/scrum-poker/backend/logic/user_logic"
	"github.com/scrum-poker/backend/logic/vote_logic"
	"github.com/scrum-poker/backend/models"
	"log"
)

func ProcessMessage(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	switch msg.Action {
	case models.ActionTypeSubmit:
		handleSubmitVote(broadcastFunc, roomId, msg)
	case models.ActionTypeReveal:
		handleRevealVotes(broadcastFunc, roomId, msg)
	case models.ActionTypeReset:
		handleResetVotes(broadcastFunc, roomId, msg)
	case models.ActionTypeTransfer:
		handleTransferScrumMaster(broadcastFunc, roomId, msg)
	case models.ActionTypeRename:
		handleRenameUser(broadcastFunc, roomId, msg)
	case models.ActionTypeLeave:
		handleLeaveRoom(broadcastFunc, roomId, msg)
	case models.ActionTypePing:
		pongMsg := &models.Message{
			Action:  models.ActionTypePong,
			Payload: msg.Payload,
		}
		broadcastFunc(roomId, pongMsg)
	default:
		broadcastFunc(roomId, msg)
	}
}

func handleSubmitVote(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for submit vote")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in submit vote payload")
		return
	}

	voteValue, ok := payload["vote"].(string)
	if !ok {
		log.Printf("Invalid vote value")
		return
	}

	err := vote_logic.SubmitVote(userId, roomId, voteValue)
	if err != nil {
		log.Printf("Failed to submit vote: %v", err)
		return
	}

	submitMsg := &models.Message{
		Action: models.ActionTypeSubmit,
		Payload: map[string]interface{}{
			"userId": userId,
			"vote":   voteValue,
		},
	}
	broadcastFunc(roomId, submitMsg)
}

func handleRevealVotes(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for reveal votes")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in reveal votes payload")
		return
	}

	room, err := db.GetRoom(roomId)
	if err != nil {
		log.Printf("Room not found: %v", err)
		return
	}

	if room.ScrumMaster != userId {
		log.Printf("Only the Scrum Master can reveal votes")
		return
	}

	revealMsg := &models.Message{
		Action: models.ActionTypeReveal,
		Payload: map[string]interface{}{
			"votes": room.Votes,
		},
	}
	broadcastFunc(roomId, revealMsg)
}

func handleResetVotes(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for reset votes")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in reset votes payload")
		return
	}
	if err := vote_logic.ResetVotes(userId, roomId); err != nil {
		log.Printf("Failed to reset votes: %v", err)
		return
	}
	broadcastFunc(roomId, msg)
}

func handleTransferScrumMaster(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for transfer scrum master")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in transfer scrum master payload")
		return
	}

	newScrumMasterId, ok := payload["newScrumMasterId"].(string)
	if !ok || newScrumMasterId == "" {
		log.Printf("Invalid or missing newScrumMasterId in transfer scrum master payload")
		return
	}
	if err := room_logic.TransferScrumMaster(userId, roomId, newScrumMasterId); err != nil {
		log.Printf("Failed to transfer scrum master: %v", err)
		return
	}

	broadcastFunc(roomId, msg)
}

func handleRenameUser(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for rename user")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in rename user payload")
		return
	}

	name, ok := payload["name"].(string)
	if !ok || name == "" {
		log.Printf("Invalid or missing name in rename user payload")
		return
	}
	if err := user_logic.RenameUser(userId, roomId, name); err != nil {
		log.Printf("Failed to rename user: %v", err)
		return
	}
	broadcastFunc(roomId, msg)
}

func handleLeaveRoom(broadcastFunc models.BroadcastFunc, roomId string, msg *models.Message) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for leave room")
		return
	}

	userId, ok := payload["userId"].(string)
	if !ok || userId == "" {
		log.Printf("Invalid or missing userId in leave room payload")
		return
	}

	if err := room_logic.LeaveRoom(roomId, userId, broadcastFunc); err != nil {
		log.Printf("Failed to leave room: %v", err)
		return
	}

	broadcastFunc(roomId, msg)
}
