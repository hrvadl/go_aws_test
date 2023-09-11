package services

import (
	"os"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var emailAttributeName = "email"

type AuthInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth interface {
	SignUp(in *AuthInput) error
	Login(in *AuthInput)
}

type AuthService struct {
	cognito  *cognitoidentityprovider.CognitoIdentityProvider
	clientID string
}

func NewAuthService(cognito *cognitoidentityprovider.CognitoIdentityProvider) *AuthService {
	return &AuthService{
		clientID: os.Getenv("COGNITO_CLIENT_ID"),
		cognito:  cognito,
	}
}

func (a *AuthService) SignUp(in *AuthInput) error {
	_, err := a.cognito.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: &in.Username,
		Password: &in.Password,
		ClientId: &a.clientID,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			a.getEmailAttribute(&in.Email),
		},
	})

	return err
}

func (a *AuthService) Login(in *AuthInput) {}

func (a *AuthService) getEmailAttribute(email *string) *cognitoidentityprovider.AttributeType {
	return &cognitoidentityprovider.AttributeType{
		Name:  &emailAttributeName,
		Value: email,
	}
}
