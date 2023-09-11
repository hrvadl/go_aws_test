package config

import "os"

type Env struct {
	AWSRegion       string `json:"aws_region"`
	CognitoID       string `json:"cognito_id"`
	UserPoolID      string `json:"user_pool_id"`
	JWTKey          string `json:"jwt_key"`
	CognitoClientId string `json:"cognito_client_id"`
}

func NewEnv() *Env {
	return &Env{
		CognitoID:       os.Getenv("COGNITO_CLIENT_ID"),
		JWTKey:          os.Getenv("JWT_KEY"),
		UserPoolID:      os.Getenv("USER_POOL_ID"),
		AWSRegion:       os.Getenv("AWS_REGION"),
		CognitoClientId: os.Getenv("COGNITO_CLIENT_ID"),
	}
}
