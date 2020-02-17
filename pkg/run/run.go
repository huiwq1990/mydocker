package run

import (
	"github.com/huiwq1990/mydocker/pkg/cgroups"
	"github.com/huiwq1990/mydocker/pkg/container"
	"github.com/huiwq1990/mydocker/pkg/network"
	"github.com/huiwq1990/mydocker/pkg/types"
	"github.com/huiwq1990/mydocker/pkg/volume/overlay"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Run(config types.RunCommandConfig,res *types.ResourceConfig) error {

	containerID := randStringBytes(10)
	if config.Name == "" {
		config.Name = containerID
	}
	log.Infof("run action, container name:%s",config.Name)

	parentCmd, writePipe,err := NewParentProcess(config.Tty, config.Name, config.Volume, config.ImageName, nil)
	if err != nil {
		return err
	}

	// 运行init命令，启动容器
	if err := parentCmd.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	err = container.RecordContainer(parentCmd, config.CmdArray, config.Name, containerID, config.Volume)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return err
	}

	// 设置进程的cgroup
	log.Debugf("run command, start cgroup set.")
	cgroupManager := cgroups.NewCgroupManager(containerID)
	defer cgroupManager.Destroy()
	cgroupManager.SetConfig(res)
	cgroupManager.Apply(parentCmd.Process.Pid)

	// 设置容器网络
	if config.Net != "" {
		c,_ := container.GetContainer(containerID)
		if err := network.Connect(config.Net, c); err != nil {
			return err
		}
	}

	sendInitCommand(config.CmdArray, writePipe)

	if config.Tty {
		parentCmd.Wait()
		log.Debugf("container init process exit.")
		container.DeleteContainer(config.Name)
		overlay.DeleteWorkSpace(config.Volume, config.Name)
	}

	return nil
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("run action, write action: %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
