package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() *session.Session {
	return session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
}
