package service

import (
	"IM_BE/repository"
	"IM_BE/ws"
)

type WsService struct {
	repo repository.WsRepo
}

func NewWsService(repo repository.WsRepo) *WsService {
	return &WsService{repo: repo}
}

func (w *WsService) AddClient(client *ws.Client) {
	manager := ws.GetWsManager()

	// TODO：并发安全，底层原理
	manager.Register <- client

	// TODO：协程的内存泄露
	go client.WritePump()
	go client.ReadPump()
}
