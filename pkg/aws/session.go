package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hrvadl/go_aws_test/pkg/config"
)

func NewSession(cfg *config.Env) *session.Session {
	return session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(cfg.AWSRegion),
		}))
}
