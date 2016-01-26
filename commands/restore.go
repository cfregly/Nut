package commands

import (
	"errors"
	"flag"
	"fmt"
	"github.com/PagerDuty/nut/specification"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"strings"
)

type RestoreCommand struct{}

func Restore() (cli.Command, error) {
	command := &RestoreCommand{}
	return command, nil
}

func (command *RestoreCommand) Help() string {
	helpText := `
	Usage: nut restore [options] <container> <image>

	nut restore is used to create container from archived images

	-sudo    Use sudo while invoking tar
	`
	return strings.TrimSpace(helpText)
}

func (command *RestoreCommand) Synopsis() string {
	return "Create container from tarball image"
}

func (command *RestoreCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("restore", flag.ExitOnError)
	flagSet.Usage = func() { fmt.Println(command.Help()) }
	sudo := flagSet.Bool("sudo", false, "Use sudo while invoking tar")
	if err := flagSet.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}

	args = flagSet.Args()
	if len(args) != 2 {
		log.Errorln(errors.New("Insufficient argument. Please pass container name and image file name"))
		return -1
	}

	if err := specification.DecompressImage(args[0], args[1], *sudo); err != nil {
		log.Errorf("Failed to restore container. Error: %s\n", err)
		return -1
	}

	if err := specification.UpdateUTS(args[0]); err != nil {
		log.Errorf("Failed to restore container. Error: %s\n", err)
		return -1
	}
	return 0
}
