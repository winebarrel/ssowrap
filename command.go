package ssowrap

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

type Command []string

func (command Command) String() string {
	return strings.Join(command, " ")
}

func (command Command) Run(ctx context.Context, creds *Credentials) error {
	name := command[0]
	args := []string{}

	if len(command) >= 2 {
		args = command[1:]
	}

	env := os.Environ()
	credsEnv, err := creds.EnvSet()

	if err != nil {
		panic(err)
	}

	for name, value := range credsEnv {
		env = append(env, name+"="+value)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	return cmd.Run()
}
