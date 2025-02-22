package ssowrap

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Options struct {
	Command        Command    `arg:"" required:"" help:"Command and arguments."`
	AWSProfile     string     `required:"" env:"AWS_PROFILE" help:"AWS CLI profile."`
	AWSConfigFile  string     `default:"~/.aws/config" env:"AWS_CONFIG_FILE" help:"AWS CLI config file location."`
	AWSSSOCacheDir string     `default:"~/.aws/sso/cache" env:"AWS_SSO_CACHE_DIR" help:"AWS SSO token cache dir location."`
	awsConfig      aws.Config `kong:"-"`
}

func (options *Options) AfterApply() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	if strings.HasPrefix(options.AWSConfigFile, "~/") {
		options.AWSConfigFile = strings.Replace(options.AWSConfigFile, "~", homeDir, 1)
	}

	if strings.HasPrefix(options.AWSSSOCacheDir, "~/") {
		options.AWSSSOCacheDir = strings.Replace(options.AWSSSOCacheDir, "~", homeDir, 1)
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	options.awsConfig = awsCfg

	return nil
}
