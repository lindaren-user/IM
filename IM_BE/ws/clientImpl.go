package ws

import (
	"IM_BE/dto"
	manager2 "IM_BE/manager"
	"IM_BE/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ClientImpl struct {
	Id      uint64
	Conn    *websocket.Conn
	Message chan dto.MessageDto
	Manager *manager2.WsManager
}

func NewClient(id uint64, conn *websocket.Conn) *ClientImpl {
	return &ClientImpl{
		Id:      id,
		Conn:    conn,
		Message: make(chan dto.MessageDto),
		Manager: manager2.GetManager(),
	}
}

func (c *ClientImpl) GetId() uint64 {
	return c.Id
}

func (c *ClientImpl) GetMessage(message dto.MessageDto) {
	c.Message <- message
}

func (c *ClientImpl) WritePump() {
	utils.GetLogger().Debug("开启 WritePump 协程")

	for message := range c.Message {
		message.SenderId = c.Id
		message.ToId = 0

		messageBytes, err := json.Marshal(message)
		if err != nil {
			utils.GetLogger().Error("解析消息体失败", zap.Error(err))
			continue
		}

		if err = c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			utils.GetLogger().Error(fmt.Sprintf("消息发送给 %d 失败", c.Id), zap.Error(err))
			return
		}
		utils.GetLogger().Debug("发送消息成功", zap.String("message", string(messageBytes)))
	}

	//// 手动写法
	//for {
	//	message, ok := <-c.Message
	//	if !ok {
	//		break // 管道已关闭
	//	}
	//	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
	//		fmt.Println(fmt.Sprintf("消息发送给 %d 失败", c.Id))
	//	}
	//}
}

func (c *ClientImpl) ReadPump() {
	utils.GetLogger().Debug("开启 ReadPump 协程")
	defer func() {
		if err := c.Conn.Close(); err != nil {
			utils.GetLogger().Error("ws 断开失败", zap.Error(err))
		}
		utils.GetLogger().Info("ws 断开成功")

		c.Manager.Unregister <- c
	}()

	for {
		contentType, messageBytes, err := c.Conn.ReadMessage() // 一直等待来自 WebSocket 的下一条消息。知道连接断开
		if err != nil {
			utils.GetLogger().Error("读取消息失败", zap.Error(err))
			break
			// TODO:return
		}
		if contentType == websocket.TextMessage {
			var message dto.MessageDto
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				utils.GetLogger().Error("解析消息体失败", zap.Error(err))
				continue
			}
			utils.GetLogger().Debug("获取消息成功", zap.String("message", string(messageBytes)))
			c.Manager.MessageHub <- message
		}
	}
}
