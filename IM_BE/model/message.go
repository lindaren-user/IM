package model

import (
	"time"
)

type MessageType string

const (
	TextMessage   MessageType = "text"
	ImageMessage  MessageType = "image"
	FileMessage   MessageType = "file"
	SystemNotice  MessageType = "system"
	FriendRequest MessageType = "friendRequest"
	AudioMessage  MessageType = "audio"
)

type PrivateMessage struct {
	Id          uint64
	SenderId    uint64
	ReceiverId  uint64
	ContentType MessageType
	Content     string
	CreatedAt   time.Time
}

type GroupMessage struct {
	Id          uint64
	SenderId    uint64
	GroupId     uint64
	ContentType MessageType
	Content     string
	CreatedAt   time.Time
}
