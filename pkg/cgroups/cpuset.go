package cgroups

import(
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/huiwq1990/mydocker/pkg/types"
	"io/ioutil"
	"path"
	"os"
	"strconv"
)

type CpusetSubSystem struct {

}

func (s *CpusetSubSystem) Set(cgroupPath string, res *types.ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}

	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}
	return nil
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (s *CpusetSubSystem)Apply(cgroupPath string, pid int) error {

	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}

	fn := path.Join(subsysCgroupPath, "tasks")

	logrus.Debugf("cpuset write to file:%s, processId:%s",fn,pid)
	if err := ioutil.WriteFile(fn,  []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}


func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}
