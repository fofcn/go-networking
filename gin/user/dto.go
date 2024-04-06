package user

type LoginDto struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}
