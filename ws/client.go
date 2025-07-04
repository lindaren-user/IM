package ws

import (
	"IM_BE/db/redis"
	"IM_BE/dto"
	"IM_BE/mq"
	"IM_BE/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Client 基础通信结构
type Client struct {
	id                uint64
	conn              *websocket.Conn
	messageBox        chan *dto.MessageRespDto // 信箱
	Manager           *WsManager               // 有点像 gin 的 RouterGroup 的属性 Engine，指向所属的 engine
	streamSubscribers []*mq.StreamSubscriber
}

func NewClient(id uint64, conn *websocket.Conn) *Client {
	return &Client{
		id:                id,
		conn:              conn,
		messageBox:        make(chan *dto.MessageRespDto),
		Manager:           GetWsManager(),
		streamSubscribers: []*mq.StreamSubscriber{},
	}
}

func (c *Client) GetId() uint64 { return c.id }

func (c *Client) GetMessage(message *dto.MessageRespDto) { c.messageBox <- message }

func (c *Client) AddSubscriber(subscriber *mq.StreamSubscriber) {
	c.streamSubscribers = append(c.streamSubscribers, subscriber)
}

func (c *Client) closeSubscribers() {
	for _, sub := range c.streamSubscribers {
		sub.Cancel()
	}
}

func (c *Client) WritePump() {
	utils.GetLogger().Debug("开启 WritePump 协程")

	for message := range c.messageBox {
		message.SenderId = c.id
		message.ToId = 0

		utils.GetLogger().Debug("发送消息成功111111")

		messageBytes, err := json.Marshal(message)
		if err != nil {
			utils.GetLogger().Error("解析消息体失败", zap.Error(err))
			continue
		}

		if err = c.conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			utils.GetLogger().Error(fmt.Sprintf("消息发送给 %d 失败", c.id), zap.Error(err))
			return
		}
		utils.GetLogger().Debug("发送消息成功", zap.String("messageBox", string(messageBytes)))
	}
}

// TODO:读消息--》消息持久化--》写消息，好处？
// 消息是否发送成功？？？

func (c *Client) ReadPump() {
	utils.GetLogger().Debug("开启 ReadPump 协程")
	defer func() {
		if err := c.conn.Close(); err != nil {
			utils.GetLogger().Error("ws 断开失败", zap.Error(err))
		}
		utils.GetLogger().Info("ws 断开成功")

		c.closeSubscribers()

		c.Manager.Unregister <- c
	}()

	streamBroker := mq.NewStreamBroker(redis.Get())

	for {
		contentType, messageBytes, err := c.conn.ReadMessage() // 一直等待来自 WebSocket 的下一条消息。知道连接断开
		if err != nil {
			utils.GetLogger().Error("读取消息失败", zap.Error(err))
			break
		}

		var message *dto.MessageRespDto
		if contentType == websocket.TextMessage {
			message, err = utils.MessageChanger(messageBytes)
			if err != nil {
				break
			}

			message.SenderId = c.id
		}
		// ...一系列的消息处理

		var stream string
		if message != nil {
			if message.ChatType == dto.PrivateChat {
				user1, user2 := c.id, message.ToId
				if user2 < user1 {
					user1, user2 = user2, user1
				}
				stream = fmt.Sprintf("chat:user:%d:%d", user1, user2)
			} else if message.ChatType == dto.GroupChat {
				stream = fmt.Sprintf("chat:group:%d", message.ToId)
			} else {
				utils.GetLogger().Error("消息类型错误")
				break
			}
		}

		// 消息放进消息队列
		if _, err := streamBroker.Publish(context.Background(), stream, message); stream != "" && err != nil {
			utils.GetLogger().Error("消息放进 mq 失败", zap.Error(err))
			break
		}
	}
}
