package commands

import (
	"github.com/mitchellh/cli"
	"strings"
)

type RunCommand struct{}

func Run() (cli.Command, error) {
	command := &RunCommand{}
	return command, nil
}

func (command *RunCommand) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (command *RunCommand) Synopsis() string {
	return "Create tarball images of existing container"
}

func (command *RunCommand) Run(args []string) int {
	return 0
}
