package main

import (
	"github.com/PagerDuty/nut/commands"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func usage() string {
	var helpText = `
Usage: nut [options]

	Build containers using LXC runtime with pluggable build DSLs

Options:

	-help        Show usage
	-version     Print version information
	`
	return strings.TrimSpace(helpText)
}

const (
	version = "0.0.1"
)

func main() {
	c := cli.NewCLI("nut", version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"build": commands.Build,
	}
	exitStatus, err := c.Run()
	if err != nil {
		log.Errorln(err)
	}
	os.Exit(exitStatus)
}
