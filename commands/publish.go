package commands

import (
	"github.com/mitchellh/cli"
	"strings"
)

type PublishCommand struct{}

func Publish() (cli.Command, error) {
	command := &PublishCommand{}
	return command, nil
}

func (command *PublishCommand) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (command *PublishCommand) Synopsis() string {
	return "Create tarball images of existing container"
}

func (command *PublishCommand) Run(args []string) int {
	return 0
}
