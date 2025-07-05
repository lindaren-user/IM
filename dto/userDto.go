package dto

type UserLoginReqDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO:前端一定需要userId吗，不是不安全吗？？？

type UserDTO struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type FriendDTO struct {
	Id       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
