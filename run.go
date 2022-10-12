package main

import (
	cgroups "linc/Cgroups"
	"os"
	"strings"

	"linc/Cgroups/subsystems"
	"linc/container"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig,volume string) {
	parent, writePipe := container.NewParentProcess(tty,volume)
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
	//交互式创建容器，父进程等待子进程结束；默认不等待，由ID为1的init进程接管容器进程
	if tty {
		parent.Wait()
	}
	
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	//init进程结束后
	container.DeleteWorkSpace(rootURL, mntURL,volume)
	//os.Exit(0)不该存在，否则cgroupManager.Destroy()不会执行
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
