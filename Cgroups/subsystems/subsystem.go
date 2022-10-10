package subsystems

type ResourceConfig struct{
	MemoryLimit string
	CpuShare	string
	CpuSet		string
}

type Subsystem interface {
	//返回subsystem的类型
	Name() string
	//
	Set(path string, res *ResourceConfig) error
	//
	Apply(path string, pid int) error
	//
	Remove(path string) error
}

//subsystem初始化实例
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
