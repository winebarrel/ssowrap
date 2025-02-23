# ssowrap

[![CI](https://github.com/winebarrel/ssowrap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/ssowrap/actions/workflows/ci.yml)

ssowrap is a single binary tool that run a command using AWS SSO credentials.

## Installation

```sh
brew install winebarrel/ssowrap/ssowrap
```

## Usage

```
Usage: ssowrap --aws-profile=STRING <command> [flags]

Arguments:
  <command>    Command and arguments.

Flags:
  -h, --help                                   Show help.
      --aws-profile=STRING                     AWS CLI profile ($AWS_PROFILE).
      --aws-config-file="~/.aws/config"        AWS CLI config file location ($AWS_CONFIG_FILE).
      --awssso-cache-dir="~/.aws/sso/cache"    AWS SSO token cache dir location ($AWS_SSO_CACHE_DIR).
      --version
```

```sh
export AWS_PROFILE=my-profile
aws sso login
ssowrap env | grep ^AWS_
# Call AWS IAM API using curl with AWS SSO credentials.
ssowrap -- sh -c 'curl -L "https://iam.amazonaws.com/?Action=ListUsers&Version=2010-05-08" --aws-sigv4 "aws:amz:us-east-1:iam" --user "$AWS_ACCESS_KEY_ID:$AWS_SECRET_ACCESS_KEY" -H "X-Amz-Security-Token: $AWS_SESSION_TOKEN"'
```
