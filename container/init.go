package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunInit() error {
	subCmdArray := readUserCommand()
	if subCmdArray == nil || len(subCmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, subCmdArray is nil")
	}

	setUpMount()

	//在系统的PATH里找执行文件的绝对路径
	path, err := exec.LookPath(subCmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error: %v", err)
		return err
	}
	log.Infof("Find path %s ", path)

	//运行用户指定程序
	//syscall.Exec调用了Kernel的execve函数，执行对应程序即将要运行的进程，它会覆盖掉当前进程的数据、堆栈等信息，包括PID。
	//用户指定的进程运行起来，把最初的init进程替换掉
	if err := syscall.Exec(path, subCmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error: %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

/**
Init 挂载点
*/

func setUpMount() {

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	//mount proc
	//为什么要在容器中挂载/proc呢， 主要原因是因为ps、top等命令依赖于/proc目录。
	//当隔离PID的时候，ps、top等命令还是未隔离的时候一样输出。 为了让隔离空间ps、top等命令只输出当前隔离空间的进程信息。需要单独挂载/proc目录。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
