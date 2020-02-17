package container

import (
	"encoding/json"
	"github.com/huiwq1990/mydocker/pkg/types"
	"github.com/huiwq1990/mydocker/pkg/volume/overlay"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"syscall"
)

func stopContainer(containerName string) error {
	c, err := GetContainer(containerName)
	if err != nil {
		return err
	}
	pidInt, err := strconv.Atoi(c.Pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return err
	}
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerName, err)
		return err
	}

	c.Status = types.ContainerStop
	c.Pid = " "
	newContentBytes, err := json.Marshal(c)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerName, err)
		return err
	}
	dirUrl := path.Join(types.ContainerInfoLocation, containerName)
	fileUrl := path.Join(dirUrl,types.ContainerConfigName)
	if err := ioutil.WriteFile(fileUrl, newContentBytes, 0622); err != nil {
		log.Errorf("Write file %s error", fileUrl, err)

return err}
	return nil
}


func removeContainer(containerName string) {
	containerInfo, err := GetContainer(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != types.ContainerStop  {
		log.Errorf("Couldn't remove running container")
		return
	}
	dirURL := path.Join(types.ContainerInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove file %s error %v", dirURL, err)
		return
	}
	overlay.DeleteWorkSpace(containerInfo.Volume, containerName)
}
