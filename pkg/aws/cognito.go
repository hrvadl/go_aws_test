package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoIdentity interface {
}

type CognitoService struct {
}

func NewCognito(session *session.Session) *cognitoidentityprovider.CognitoIdentityProvider {
	return cognitoidentityprovider.New(session)
}
