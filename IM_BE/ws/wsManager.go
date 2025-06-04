package ws

import (
	"IM_BE/dto"
	"IM_BE/iface"
	"IM_BE/utils"
	"go.uber.org/zap"
)

// TODO:使用 redis 实现分布式
// redis 统一处理各个服务器的 WsManager

// WsManager 中转站，全局管理者，消息总线
// 方便拓展
type WsManager struct {
	Clients map[uint64]*Client
	Groups  map[uint64]map[uint64]*Client

	Register   chan *Client
	Unregister chan *Client

	MessageHub chan *dto.MessageRespDto // 读进来的消息放在这个 MessageHub 中转站

	service iface.MessageHandler
}

// TODO:sync.Once和单例模式
//// 全局单例模式
//var (
//	managerInstance *WsManager
//	onceInit        sync.Once
//)
//
//func GetManager() *WsManager {
//	onceInit.Do(func() {
//		managerInstance = newWsManager()
//		go managerInstance.run() // 开协程监听通道
//	})
//	return managerInstance
//}

var managerInstance *WsManager

func InitWsManager(handler iface.MessageHandler) {
	managerInstance = &WsManager{
		Clients:    make(map[uint64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		MessageHub: make(chan *dto.MessageRespDto),
		service:    handler,
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

		case message := <-w.MessageHub:
			// 持久化到数据库
			if err := w.service.SaveMessage(message); err != nil {
				continue
			}

			if _, ok := w.Clients[message.ToId]; ok {
				w.Clients[message.ToId].GetMessage(message)
			}
		}
	}
}
