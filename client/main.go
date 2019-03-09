package main

import (
	"codex.tunnel/commands"
	"github.com/jessevdk/go-flags"
	"os"
)

var opts struct {
	Run commands.RunCommand `command:"run" description:"Run tunnel"`
	Generate commands.GenerateCommand `command:"generate" description:"Generate rsa keys"`
	Tunnel commands.TunnelCommand `command:"tunnel" description:"Setup tunnel manually"`
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(0)
	}
}