package specification

import (
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"gopkg.in/lxc/go-lxc.v2"
	"time"
)

var (
	MinimalEnv = []string{
		"SHELL=/bin/bash",
		"USER=root",
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"PWD=/root",
		"EDITOR=vim",
		"LANG=en_US.UTF-8",
		"HOME=/root",
		"LANGUAGE=en_US",
		"LOGNAME=root",
	}
)

func UUID() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}
	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F
	return hex.EncodeToString(u), nil
}

func CloneAndStartContainer(original, cloned string) (*lxc.Container, error) {
	orig, err := lxc.NewContainer(original)
	if err != nil {
		log.Errorf("Failed to initialize container object. Error: %v", err)
		return nil, err
	}
	if err := orig.Clone(cloned, lxc.CloneOptions{}); err != nil {
		log.Errorf("Failed to clone container %s as %s. Error: %v", original, cloned, err)
		return nil, err
	}
	ct, err := lxc.NewContainer(cloned)
	if err != nil {
		log.Errorf("Failed to clone container %s as %s. Error: %v", original, cloned, err)
	}
	if err := ct.Start(); err != nil {
		log.Errorf("Failed to start cloned container %s. Error: %v", cloned, err)
		return nil, err
	}
	log.Infof("Created container named: %s. Waiting for ip allocation", cloned)
	if _, err := ct.WaitIPAddresses(30 * time.Second); err != nil {
		log.Errorf("Failed to while waiting to start the container %s. Error: %v", cloned, err)
		return nil, err
	}
	return ct, nil
}
