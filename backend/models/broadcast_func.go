package models

type BroadcastFunc func(roomId string, msg *Message)
type ConnectionChecker func(roomId, userId string) bool
