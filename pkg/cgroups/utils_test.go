package cgroups

import(
	"errors"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestE(t *testing.T){
	logrus.Errorf("aa %s %s","a",errors.New("bb"))
}

func TestFindCgroupMountpoint(t *testing.T) {
	t.Logf("cpu subsystem mount point %v\n", FindCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem mount point %v\n", FindCgroupMountpoint("cpuset"))
	t.Logf("memory subsystem mount point %v\n", FindCgroupMountpoint("memory"))
}