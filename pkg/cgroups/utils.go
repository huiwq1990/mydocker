package cgroups

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"os"
	"path"
	"bufio"
)

// /sys/fs/cgroup/memory
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				log.Debugf("find cgroup:%s, mount point:%s",subsystem,fields[4])
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("find cgroup:%s error:%s",subsystem,err)
		return ""
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	//根据subsystem需要找到当前容器的cgroup位置,这样才可以往里面加入相关的限制.
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}