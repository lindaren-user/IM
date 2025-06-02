package service

import "IM_BE/repository"

type WsService struct {
	repo repository.WsRepo
}

func NewWsService(repo repository.WsRepo) *WsService {
	return &WsService{repo: repo}
}
