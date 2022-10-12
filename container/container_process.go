package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
	
)
var(
	DefaultInfoLocation string = "/var/run/linc/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile	string = "container.log"
)

func NewParentProcess(tty bool, volume,containerName string) (*exec.Cmd, *os.File) {
	rp, wp, err := NewPipe()
	if err != nil {
		log.Errorf("NewPipe error %v", err)
		return nil, nil
	}
	//  /proc/self/exe 链接到当前进程的执行命令文件，这里是linc文件
	cmd := exec.Command("/proc/self/exe", "init")
	//clone参数就是去fork出来一个新进程
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {//logs 命令实现
		dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
			return nil, nil
		}
		stdLogFilePath := dirURL + ContainerLogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	//传入管道文件读取端的句柄；一个进程默认会有三个文件描述符，分别是Stdin、Stdout、Stderr，所以rp会成为第四个
	cmd.ExtraFiles = []*os.File{rp}

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	//init进程执行前
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL

	return cmd, wp
}

func NewPipe() (*os.File, *os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return r, w, err
}
