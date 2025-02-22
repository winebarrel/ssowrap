package ssowrap

import (
	"context"
)

func Run(ctx context.Context, options *Options) error {
	sso := NewSSO(options)
	creds, err := sso.GetCredentials(ctx)

	if err != nil {
		return err
	}

	return options.Command.Run(ctx, creds)
}
