package main

import (
	cgroups "linc/Cgroups"
	"os"
	"strings"

	"linc/Cgroups/subsystems"
	"linc/container"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("linc-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	//init进程结束后
	container.DeleteWorkSpace(rootURL, mntURL)
	//os.Exit(0)不该存在，否则cgroupManager.Destroy()不会执行
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
