package jwtauth

type TokenInfo struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	TokenType    string `json:"token_type"`    // 令牌类型
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}
