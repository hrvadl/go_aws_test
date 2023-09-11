package config

import "os"

type Env struct {
	CognitoID  string `json:"cognito_id"`
	UserPoolID string `json:"user_pool_id"`
	JWTKey     string `json:"jwt_key"`
}

func NewEnv() *Env {
	return &Env{
		CognitoID:  os.Getenv("COGNITO_CLIENT_ID"),
		JWTKey:     os.Getenv("JWT_KEY"),
		UserPoolID: os.Getenv("USER_POOL_ID"),
	}
}
