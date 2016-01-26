package main

import (
	"github.com/PagerDuty/nut/commands"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	version = "0.0.1"
)

func main() {
	c := cli.NewCLI("nut", version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"build": commands.Build,
		"fetch": commands.Fetch,
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	exitStatus, err := c.Run()
	if err != nil {
		log.Errorln(err)
	}
	os.Exit(exitStatus)
}
