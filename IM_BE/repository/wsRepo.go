package repository

import "database/sql"

type WsRepo interface {
}

func NewWsRepo(db *sql.DB) WsRepo {
	return &WsRepoImpl{db: db}
}

type WsRepoImpl struct {
	db *sql.DB
}
