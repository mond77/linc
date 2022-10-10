package container

import (
	"os"
	"syscall"
	"github.com/sirupsen/logrus"
)

func RunInit(command string, args []string) error {
	logrus.Infof("command %s", command)
	//初始化容器内容
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//mount proc文件系统
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	//运行用户指定程序
	//syscall.Exec调用了Kernel的execve函数，执行对应程序即将要运行的进程，它会覆盖掉当前进程的数据、堆栈等信息，包括PID。
	//用户指定的进程运行起来，把最初的init进程替换掉
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
