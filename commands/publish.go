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

type PublishCommand struct{}

func Publish() (cli.Command, error) {
	command := &PublishCommand{}
	return command, nil
}

func (command *PublishCommand) Help() string {
	helpText := `
	Usage: nut publish <rootfs.tgz> <region> <bucket> <key>

	nut publish is used to publish a rootfs image into s3
	Use environment variables or .credential file to pass
	aws credentials
	`
	return strings.TrimSpace(helpText)
}

func (command *PublishCommand) Synopsis() string {
	return "Publish tarball images of existing container in s3"
}

func (command *PublishCommand) Run(args []string) int {

	flagSet := flag.NewFlagSet("publish", flag.ExitOnError)
	flagSet.Usage = func() { fmt.Println(command.Help()) }
	if err := flagSet.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}

	args = flagSet.Args()
	if len(args) != 4 {
		log.Errorln(errors.New("Insufficient argument. Please pass container image file, s3 region, bucket and key"))
		return -1
	}
	if err := specification.Publish(args[0], args[1], args[2], args[3]); err != nil {
		log.Errorln(err)
		return -1
	}
	return 0
}
