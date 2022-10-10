package subsystems

import (
	"os/exec"
	"testing"
)

func TestFindCgroupMountpoint(t *testing.T) {
	cmd := exec.Command("/media/psf/Home/Documents/goProjects/linc/linc",[]string{"run","-it","sh"}...)
	cmd.Start()
	t.Logf("cpu subsystem mount point %v\n", FindCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", FindCgroupMountpoint("cpuset"))
	t.Logf("memory subsystem mount point %v\n", FindCgroupMountpoint("memory"))
	cmd.Wait()
}