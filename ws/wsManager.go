package ws

import (
	"IM_BE/utils"
	"go.uber.org/zap"
)

// TODO:使用 redis， 可以直接做到分布式
// redis 统一处理各个服务器的 WsManager

// TODO；使用 stream 订阅的方式处理聊天，群聊私聊分开，方便做不同的优化
// redis List 的不足：
// 1. 不适用于多端的信息接收，一个端接收后，另一个端无法接收（单一消费）
// 2. 对于一个群聊，其中一个人获取消息之后，其他人无法得到该消息（单一消费）
// 3. 没有 ack 机制

// redis Pub/Sub 的不足：
// 1. 没有 ack 机制
// 2. 不能持久化数据（实际上这点对于 IM 没有影响，毕竟消息是存储在数据库的）

// 推荐使用 redis Stream：
// 1. 存在 ack 机制
// 2. 有消费者组（多消费者）

// 相比于使用 channel：
// 1. 天然的分布式
// 2. 很容易的加群加好友
// 3. 私聊和群聊的消息传递方式几乎一样

// TODO：从数据库中加载群聊

// WsManager 全局管理者，方便拓展
type WsManager struct {
	Clients map[uint64]*Client

	Register   chan *Client
	Unregister chan *Client
}

var managerInstance *WsManager

func InitWsManager() {
	managerInstance = &WsManager{
		Clients:    make(map[uint64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	go managerInstance.run()
}

func GetWsManager() *WsManager {
	return managerInstance
}

func (w *WsManager) run() {
	for {
		select {
		case client := <-w.Register:
			w.Clients[client.GetId()] = client
			utils.GetLogger().Info("用户连接", zap.Uint64("userId", client.GetId()))

		case client := <-w.Unregister:
			if _, ok := w.Clients[client.GetId()]; ok {
				delete(w.Clients, client.GetId())
				utils.GetLogger().Info("用户断开连接", zap.Uint64("userId", client.GetId()))
			}
		}
	}
}
