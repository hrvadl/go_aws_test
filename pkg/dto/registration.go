package dto

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

var emailAttributeName = "email"

type SignUpDTO struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (o *SignUpDTO) GetCognitoEmailAttribute() *cognitoidentityprovider.AttributeType {
	return &cognitoidentityprovider.AttributeType{
		Name:  &emailAttributeName,
		Value: &o.Email,
	}
}
