package user

type LoginDto struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type UserInfoDto struct {
	Username string `json:"username"`
	UserId   uint   `json:"userId"`
}
