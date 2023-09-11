package dto

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

var emailAttributeName = "email"
var emailVerifiedAttributeName = "email_verified"

type SignInDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthDTO(email, username, password string) *SignInDTO {
	return &SignInDTO{
		Email:    email,
		Username: username,
		Password: password,
	}
}

func (o *SignInDTO) GetCognitoEmailAttribute() *cognitoidentityprovider.AttributeType {
	return &cognitoidentityprovider.AttributeType{
		Name:  &emailAttributeName,
		Value: &o.Email,
	}
}

// For the sake of simplicity skip email verification
func (o *SignInDTO) GetCognitoEmailVerifiedAttribute() *cognitoidentityprovider.AttributeType {
	verified := "true"
	return &cognitoidentityprovider.AttributeType{
		Name:  &emailVerifiedAttributeName,
		Value: &verified,
	}
}
