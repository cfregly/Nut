package main

import (
	"github.com/PagerDuty/nut/specification"
	log "github.com/sirupsen/logrus"
)

func main() {
	file := "Dockerfile"
	name, err := specification.UUID()
	if err != nil {
		log.Fatalf("Failed to create uuid. Error: %s\n", err)
	}
	spec := specification.New(name, file)
	if err := spec.Parse(); err != nil {
		log.Fatalf("Failed to parse dockerfile. Error: %s\n", err)
	}
	if err := spec.Build(); err != nil {
		log.Fatalf("Failed to build container from dockerfile. Error: %s\n", err)
	}
}
