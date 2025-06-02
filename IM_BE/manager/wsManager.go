package manager

import (
	"IM_BE/dto"
	"IM_BE/iface"
	"IM_BE/utils"
	"go.uber.org/zap"
	"sync"
)

// WsManager 中转站，全局管理者，消息总线
// 方便拓展
type WsManager struct {
	Clients map[uint64]iface.Client

	Groups map[uint64]map[uint64]iface.Client

	Register chan iface.Client

	Unregister chan iface.Client

	MessageHub chan dto.MessageDto
}

// TODO:sync.Once和单例模式
// 全局单例模式
var (
	managerInstance *WsManager
	onceInit        sync.Once
)

func GetManager() *WsManager {
	onceInit.Do(func() {
		managerInstance = newWsManager()
		go managerInstance.run() // 开协程监听通道
	})
	return managerInstance
}

func newWsManager() *WsManager {
	return &WsManager{
		Clients:    make(map[uint64]iface.Client),
		Register:   make(chan iface.Client),
		Unregister: make(chan iface.Client),
		MessageHub: make(chan dto.MessageDto),
	}
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
			if _, ok := w.Clients[message.ToId]; ok {
				w.Clients[message.ToId].GetMessage(message)
			}
		}
	}
}
