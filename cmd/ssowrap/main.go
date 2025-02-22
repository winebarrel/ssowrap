package main

import (
	"context"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/ssowrap"
)

var version string

func init() {
	log.SetFlags(0)
}

func parseArgs() *ssowrap.Options {
	var cli struct {
		ssowrap.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	ctx := context.Background()
	err := ssowrap.Run(ctx, options)

	if err != nil {
		log.Fatalf("ssowrap: %s: %s", options.Command, err)
	}
}
