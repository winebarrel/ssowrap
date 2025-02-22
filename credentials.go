package ssowrap

import (
	"github.com/Netflix/go-env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
)

type Credentials struct {
	AccessKeyId     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	SessionToken    string `env:"AWS_SESSION_TOKEN"`
}

func NewCredsFromRoleCreds(roleCreds *types.RoleCredentials) *Credentials {
	creds := &Credentials{
		AccessKeyId:     aws.ToString(roleCreds.AccessKeyId),
		SecretAccessKey: aws.ToString(roleCreds.SecretAccessKey),
		SessionToken:    aws.ToString(roleCreds.SessionToken),
	}

	return creds
}

func (creds *Credentials) EnvSet() (env.EnvSet, error) {
	return env.Marshal(creds)
}
