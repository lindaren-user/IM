package model

import "time"

type User struct {
	ID       uint
	Username string
	Password string
	Nickname string
	Avatar   string
	Email    string
	Status   int
	CreateAt time.Time
	UpdateAt time.Time
}
