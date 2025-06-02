package model

import (
	"encoding/json"
	"time"
)

type PrivateMessage struct {
	SenderId   uint64
	ReceiverId uint64
	Content    json.RawMessage
	CreatedAt  time.Time
}

type GroupMessage struct {
	SenderId  uint64
	GroupId   uint64
	Content   json.RawMessage
	CreatedAt time.Time
}
