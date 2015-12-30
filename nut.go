package main

import (
	"flag"
	"fmt"
	"github.com/PagerDuty/nut/specification"
	log "github.com/sirupsen/logrus"
	"strings"
)

func usage() string {
	var helpText = `
Usage: nut [options]

	Build containers using LXC runtime with pluggable build DSLs

Options:

	-help        Show usage
	-specfile    Local path to the specification file (defaults to dockerfle)
	-ephemeral   Destroy the container after creation
	-name        Name of the container (defaults to randomly generated UUID)
	-version    Print version information
	`
	return strings.TrimSpace(helpText)
}

const (
	version = "0.0.1"
)

func main() {
	file := flag.String("specfile", "Dockerfile", "Container build specification file")
	help := flag.Bool("help", false, "Show usage")
	ephemeral := flag.Bool("ephemeral", false, "Destroy the container after creating it")
	versionFalg := flag.Bool("version", false, "Print version information")
	name := flag.String("name", "", "Name of the resulting container (defaults to randomly generated UUID)")
	flag.Parse()
	if *help {
		fmt.Println(usage())
		return
	}

	if *versionFalg {
		fmt.Println(version)
		return
	}

	if *name == "" {
		uuid, err := specification.UUID()
		if err != nil {
			log.Fatalf("Failed to create uuid. Error: %s\n", err)
		}
		name = &uuid
	}
	spec := specification.New(*name, *file)
	if err := spec.Parse(); err != nil {
		log.Fatalf("Failed to parse dockerfile. Error: %s\n", err)
	}
	if err := spec.Build(); err != nil {
		log.Fatalf("Failed to build container from dockerfile. Error: %s\n", err)
	}
	if *ephemeral {
		log.Infof("Ephemeral mode. Destroying the container")
	}
}
