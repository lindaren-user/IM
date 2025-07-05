package dto

import (
	"IM_BE/model"
	"encoding/json"
	"strconv"
	"time"
)

type ChatType string

const (
	PrivateChat ChatType = "private"
	GroupChat   ChatType = "group"
)

// TODO: copier无法复制匿名字段
type messageBasic struct {
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"`
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
	Seq         uint64            `json:"seq"`
}

type MessageReqDto struct {
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"`
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
	Content     json.RawMessage   `json:"content"`
	Seq         uint64            `json:"seq"`
}

type RawMessage struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	ToID        string `json:"to_id"`
	ChatType    string `json:"chat_type"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
	Content     string `json:"content"`
	Seq         string `json:"seq"`
}

func (r *RawMessage) ToMessageRespDto() (*MessageRespDto, error) {
	id, err := strconv.ParseUint(r.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	senderID, err := strconv.ParseUint(r.SenderID, 10, 64)
	if err != nil {
		return nil, err
	}
	toID, err := strconv.ParseUint(r.ToID, 10, 64)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	seq, err := strconv.ParseUint(r.Seq, 10, 64)
	if err != nil {
		seq = 0
	}

	return &MessageRespDto{
		Id:          id,
		SenderId:    senderID,
		ToId:        toID,
		ChatType:    ChatType(r.ChatType),
		ContentType: model.MessageType(r.ContentType),
		CreatedAt:   createdAt,
		Content:     r.Content,
		Seq:         seq,
	}, nil
}

type MessageRespDto struct {
	Id          uint64            `json:"id"`
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"` // TODO：多余的部分，等待删除
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
	Content     string            `json:"content"`
	Seq         uint64            `json:"seq"`
}

func (m *MessageRespDto) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           m.Id,
		"sender_id":    m.SenderId,
		"to_id":        m.ToId,
		"chat_type":    string(m.ChatType),
		"content_type": string(m.ContentType),
		"created_at":   m.CreatedAt.Format(time.RFC3339),
		"content":      m.Content,
		"seq":          m.Seq,
	}
}
