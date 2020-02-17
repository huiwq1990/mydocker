package container

import (
	"encoding/json"
	"github.com/huiwq1990/mydocker/pkg/types"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetContainer(cid string) (*types.ContainerInfo, error) {
	content, err := ioutil.ReadFile(path.Join(types.ContainerInfoLocation, cid, types.ContainerConfigName))
	if err != nil {
		return nil, err
	}
	var containerInfo types.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}

func RecordContainer(pproc *exec.Cmd, commandArray []string, containerName, id, volume string) (error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	dirUrl := path.Join(types.ContainerInfoLocation, containerName)
	fileUrl := path.Join(dirUrl,types.ContainerConfigName)
	command := strings.Join(commandArray, "")
	containerInfo := &types.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(pproc.Process.Pid),
		Command:     command,
		CreatedTime: createTime,
		Status:      types.ContainerRunning,
		Name:        containerName,
		Volume:      volume,
		MountUrl: pproc.Path,
		LogFile: path.Join(fileUrl,types.ContainerLogFileName),
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return err
	}
	jsonStr := string(jsonBytes)

	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return err
	}
	file, err := os.Create(fileUrl)
	log.Debugf("create container file: %s",dirUrl)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileUrl, err)
		return err
	}

	log.Debugf("write container info. %s",jsonStr)
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return err
	}

	return nil
}

func DeleteContainer(containerId string) error{
	if err := os.RemoveAll(path.Join(types.ContainerInfoLocation,containerId)); err != nil {
		return err
	}
	return nil
}

func CreateLogFile(containerName string) (*os.File,error) {
	dirUrl := path.Join(types.ContainerInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return nil,err
	}

	logFile, err := os.Create(path.Join(dirUrl,types.ContainerLogFileName))
	if err != nil {
		return nil,err
	}
	return logFile,nil
}
