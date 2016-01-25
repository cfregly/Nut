package commands

import (
	"flag"
	"github.com/PagerDuty/nut/specification"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type BuildCommand struct {
}

func Build() (cli.Command, error) {
	command := &BuildCommand{}
	return command, nil
}

func (command *BuildCommand) Help() string {
	helpText := `
		-specfile    Local path to the specification file (defaults to dockerfle)
		-ephemeral   Destroy the container after creation
		-name        Name of the container (defaults to randomly generated UUID)
		-stop        Stop container at the end
		-volume      Mount host directory inside container
		-export      Create tarball of container rootfs (assumes -stop)
		-publish     Publish tarball in S3
	`
	return strings.TrimSpace(helpText)
}

func (command *BuildCommand) Synopsis() string {
	synopsis := "Build container from Dockerfile"
	return synopsis
}

func (command *BuildCommand) Run(args []string) int {

	flagSet := flag.NewFlagSet("build", flag.ExitOnError)

	file := flagSet.String("specfile", "Dockerfile", "Container build specification file")
	stopAfterBuild := flagSet.Bool("stop", false, "Stop container after build")
	ephemeral := flagSet.Bool("ephemeral", false, "Destroy the container after creating it")
	name := flagSet.String("name", "", "Name of the resulting container (defaults to randomly generated UUID)")
	export := flagSet.String("export", "", "File path for the container tarball")
	exportSudo := flagSet.Bool("export-sudo", false, "Use sudo while invoking tar")
	volume := flagSet.String("volume", "", "Mount host directory inside container. Format: '[host_directory:]container_directory[:mount options]")
	publish := flagSet.String("publish", "", "Publish tarball in s3 (assumes -stop, -export)")

	flagSet.Parse(args)
	if *publish != "" {
		if *export == "" {
			uuid, err := specification.UUID()
			if err != nil {
				log.Errorln(err)
				return -1
			}
			*export = uuid
		}
	}

	if *export != "" {
		*stopAfterBuild = true
	}
	if *name == "" {
		uuid, err := specification.UUID()
		if err != nil {
			log.Errorln(err)
			return -1
		}
		name = &uuid
	}

	spec := specification.New(*name, *file)

	if err := spec.Parse(); err != nil {
		log.Errorf("Failed to parse dockerfile. Error: %s\n", err)
		return -1
	}

	if err := spec.Build(*volume); err != nil {
		log.Errorf("Failed to build container from dockerfile. Error: %s\n", err)
		return -1
	}

	if *stopAfterBuild {
		log.Infof("Stopping container")
		if err := spec.Stop(); err != nil {
			log.Errorf("Failed to stop container. Error: %s\n", err)
			return -1
		}
	}

	if *export != "" {
		log.Infof("Exporting container")
		if err := spec.ExportContainer(*export, *exportSudo); err != nil {
			log.Errorf("Failed to export container. Error: %s\n", err)
			return -1
		}
	}

	if *ephemeral {
		log.Infof("Ephemeral mode. Destroying the container")
		if err := spec.Destroy(); err != nil {
			log.Errorf("Failed to destroy container. Error: %s\n", err)
			return -1
		}
	}
	if *publish != "" {
		parts := strings.Split(*publish, "/")
		bucket := parts[0]
		key := strings.Join(parts[1:], "/")
		svc := s3.New(session.New(), aws.NewConfig().WithRegion("us-west-1"))
		fi, err := os.Open(*export)
		if err != nil {
			log.Error(err)
			return -1
		}
		defer fi.Close()
		params := &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   fi,
		}

		_, uploadErr := svc.PutObject(params)
		if uploadErr != nil {
			log.Error(uploadErr)
			return -1
		}
	}
	return 0
}
