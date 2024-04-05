package user

type TokenService struct {
}

type CustomClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func (ts *TokenService) GenerateToken(username string) (string, error) {
	claims := CustomClaims{
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
