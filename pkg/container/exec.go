package container

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	// TODO 网络这块暂时不搞
	//_ "github.com/huiwq1990/mydocker/nsenter"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func ExecContainer(containerName string, comArray []string) error {
	container, err := GetContainer(containerName)
	if err != nil {
		return errors.Errorf("get container: %s error %v", containerName, err)
	}

	cmdStr := strings.Join(comArray, " ")
	log.Debugf("exec cmd, container pid %s, cmd: %s", container.Pid,cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, container.Pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)
	containerEnvs,err := getEnvsByPid(container.Pid)
	if err != nil {
		return err
	}
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func getEnvsByPid(pid string) ([]string,error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil,errors.Wrapf(err,"read file %s error %v", path)
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs,nil
}
