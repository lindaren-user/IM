package repository

import "database/sql"

// 名字以业务划分
type WsRepo interface {
}

func NewWsRepo(db *sql.DB) WsRepo {
	return &WsRepoImpl{db: db}
}

type WsRepoImpl struct {
	db *sql.DB
}
