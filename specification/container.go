package specification

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/lxc/go-lxc.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func SetupBindMounts(container *lxc.Container, volume string) error {
	// bind syntax: "/tmp/xxx home/ubuntu/foo none bind,create=dir"
	// cli spec:  [host_directory:]container_directory
	// For unprivileged containers rw bind mount still does not allow writing files from within the container, due to posix ACL
	// workaround: on host:
	//   setfacl -Rm user:ranjib:rwx,default:user:ranjib:rwx,user:100000:rwx,user:101000:rwx,default:user:100000:rwx,default:user:101000:rwx /tmp/xxx
	parts := strings.Split(volume, ":")
	options := []string{"none", "bind,create=dir", "0", "0"}
	var hostDir string
	var containerDir string
	switch len(parts) {
	case 1:
		containerDir = volume
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		hostDir = dir
	case 2:
		containerDir = parts[1]
		if p, err := filepath.Abs(parts[0]); err != nil {
			return err
		} else {
			hostDir = p
		}
	case 3:
		containerDir = parts[1]
		if p, err := filepath.Abs(parts[0]); err != nil {
			return err
		} else {
			hostDir = p
		}
		options[1] = "bind," + parts[2]
	default:
		fmt.Errorf("Invalid volume spec. Parts: %d", len(parts))
	}
	containerDir = strings.TrimPrefix(containerDir, "/")
	val := hostDir + " " + containerDir + " " + strings.Join(options, " ")
	log.Debugf("Setting up bind mounts: %s\n", val)
	path := container.ConfigFileName()
	if err := container.SetConfigItem("lxc.mount.entry", val); err != nil {
		return err
	}
	if err := container.SaveConfigFile(path); err != nil {
		return err
	}
	return nil
}

func CloneAndStartContainer(original, cloned, volume string) (*lxc.Container, error) {
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
		return nil, err
	}

	if volume != "" {
		if err := SetupBindMounts(ct, volume); err != nil {
			log.Errorf("Failed to setup volumes for %s. Error: %v", cloned, err)
			return nil, err
		}
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

func ExportContainer(name, file string, sudo, pigz bool) error {
	//command := "sudo tar -Jcpf rootfs1.tar.xz -C ~/.local/share/lxc/ruby_2.3/rootfs  . --numeric-owner"
	lxcdir := lxc.GlobalConfigItem("lxc.lxcpath")
	ctDir := filepath.Join(lxcdir, name)
	var command string
	if pigz {
	} else {
		command = fmt.Sprintf("tar -Jcpf %s --numeric-owner -C %s .", file, ctDir)
	}
	if sudo {
		command = "sudo " + command
	}
	log.Infof("Invoking: %s", command)
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		log.Error(err)
		return err
	}
	return nil
}

func DecompressImage(name, file string, sudo bool) error {
	lxcpath := lxc.GlobalConfigItem("lxc.lxcpath")
	ctDir := filepath.Join(lxcpath, name)
	untarCommand := fmt.Sprintf("tar --numeric-owner -xpJf  %s -C %s", file, ctDir)
	if sudo {
		untarCommand = "sudo " + untarCommand
	}
	if err := os.Mkdir(ctDir, 0770); err != nil {
		log.Errorln(err)
		return err
	}
	log.Infof("Invoking: %s", untarCommand)
	parts := strings.Fields(untarCommand)
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		log.Error(err)
		return err
	}
	return nil
}

func UpdateUTS(name string) error {
	ct, err := lxc.NewContainer(name)
	rootfs := filepath.Join(lxc.GlobalConfigItem("lxc.lxcpath"), name, "rootfs")
	if err != nil {
		return err
	}
	if err := ct.LoadConfigFile(ct.ConfigFileName()); err != nil {
		return err
	}
	if err := ct.SetConfigItem("lxc.utsname", name); err != nil {
		return err
	}
	if err := ct.SetConfigItem("lxc.rootfs", rootfs); err != nil {
		return err
	}
	return ct.SaveConfigFile(ct.ConfigFileName())
}
