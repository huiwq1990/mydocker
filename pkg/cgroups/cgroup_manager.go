package cgroups

import (
	log "github.com/sirupsen/logrus"
	"github.com/huiwq1990/mydocker/pkg/types"
)


var (
	SubsystemsIns = []types.Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)

type CgroupManager struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	// 这里用的是containerId
	Path     string
	// 资源配置
	Resource *types.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程pid加入到这个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置cgroup资源限制
func (c *CgroupManager) SetConfig(res *types.ResourceConfig) error {
	for _, subSysIns := range(SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

//释放cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range(SubsystemsIns) {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
