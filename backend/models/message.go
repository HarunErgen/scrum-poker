package models

type ActionType string

const (
	ActionTypeJoin     ActionType = "join"
	ActionTypeOffline  ActionType = "offline"
	ActionTypeOnline   ActionType = "online"
	ActionTypeLeave    ActionType = "leave"
	ActionTypeRename   ActionType = "rename"
	ActionTypeSubmit   ActionType = "submit"
	ActionTypeReveal   ActionType = "reveal"
	ActionTypeReset    ActionType = "reset"
	ActionTypeTransfer ActionType = "transfer"
)

type Message struct {
	Action  ActionType  `json:"action"`
	Payload interface{} `json:"payload"`
}
