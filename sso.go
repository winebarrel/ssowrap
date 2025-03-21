package ssowrap

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/bigkevmcd/go-configparser"
)

type SSO struct {
	*Options
	c *sso.Client
}

func NewSSO(options *Options) *SSO {
	sso := &SSO{
		Options: options,
		c:       sso.NewFromConfig(options.awsConfig),
	}

	return sso
}

func (sso *SSO) getSSOStartUrl() (string, error) {
	cfg, err := configparser.NewConfigParserFromFile(sso.AWSConfigFile)

	if err != nil {
		return "", err
	}

	return cfg.Get("profile "+sso.AWSProfile, "sso_start_url")
}

func (ssoClient *SSO) GetCredentials(ctx context.Context) (*Credentials, error) {
	ssoStartUrl, err := ssoClient.getSSOStartUrl()

	if err != nil {
		return nil, err
	}

	ssoCache := NewSSOCache(ssoClient.Options, ssoStartUrl)
	token, err := ssoCache.LastToken()

	if err != nil {
		return nil, err
	}

	sts := NewSTS(ssoClient.Options)
	role, err := sts.GetRole(ctx)

	if err != nil {
		return nil, err
	}

	input := &sso.GetRoleCredentialsInput{
		AccountId:   aws.String(role.Account),
		RoleName:    aws.String(role.Name),
		AccessToken: aws.String(token.AccessToken),
	}

	roleCreds, err := ssoClient.c.GetRoleCredentials(ctx, input, func(o *sso.Options) {
		o.Region = token.Region
	})

	if err != nil {
		return nil, err
	}

	creds := &Credentials{
		AccessKeyId:     aws.ToString(roleCreds.RoleCredentials.AccessKeyId),
		SecretAccessKey: aws.ToString(roleCreds.RoleCredentials.SecretAccessKey),
		SessionToken:    aws.ToString(roleCreds.RoleCredentials.SessionToken),
	}

	return creds, nil
}
