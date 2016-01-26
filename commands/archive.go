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

type ArchiveCommand struct{}

func Archive() (cli.Command, error) {
	command := &ArchiveCommand{}
	return command, nil
}

func (command *ArchiveCommand) Help() string {
	helpText := `
	Usage: nut archive [options] <container> <image>

	nut archive is used to build tarball image from an existing
	container.

	-sudo    Use sudo while invoking tar
	`
	return strings.TrimSpace(helpText)
}

func (command *ArchiveCommand) Synopsis() string {
	return "Create tarball images of existing container"
}

func (command *ArchiveCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("archive", flag.ExitOnError)
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

	if err := specification.ExportContainer(args[0], args[1], *sudo); err != nil {
		log.Errorf("Failed to export container. Error: %s\n", err)
		return -1
	}
	return 0
}
