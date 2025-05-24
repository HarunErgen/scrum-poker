package models

type BroadcastFunc func(roomId string, msg *Message)
