package ssowrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STS struct {
	*Options
	c *sts.Client
}

func NewSTS(options *Options) *STS {
	sts := &STS{
		Options: options,
		c:       sts.NewFromConfig(options.awsConfig),
	}

	return sts
}

type Role struct {
	Account string
	Name    string
}

func parseCallerIdentityArn(arn string) (string, error) {
	arnParts := strings.Split(arn, "/")

	if len(arnParts) != 3 {
		return "", fmt.Errorf("failed to parse caller identity ARN: %s", arn)
	}

	role := arnParts[1]
	roleParts := strings.Split(role, "_")

	if len(roleParts) != 3 {
		return "", fmt.Errorf("failed to parse role IN ARN: %s", role)
	}

	return roleParts[1], nil
}

func (stsClient *STS) GetRole(ctx context.Context) (*Role, error) {
	output, err := stsClient.c.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	arn := aws.ToString(output.Arn)
	roleName, err := parseCallerIdentityArn(arn)

	if err != nil {
		return nil, err
	}

	role := &Role{
		Account: aws.ToString(output.Account),
		Name:    roleName,
	}

	return role, nil
}
