package run

import (
	"github.com/huiwq1990/mydocker/pkg/container"
	"github.com/huiwq1990/mydocker/pkg/volume/overlay"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// init进程作用
// 当执行 ps -ef时还是可以看到整个宿主机的进程.
// 但是先执行mount -t proc proc /proc后在执行ps -ef的时候就只显示当前namespace内进程的状态了,
func NewParentProcess(tty bool, containerName, volume, imageName string, envSlice []string) (*exec.Cmd, *os.File, error) {
	// 处理管道
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil,err
	}
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorf("get init process error %v", err)
		return nil, nil,err
	}
	log.Debugf("run cmd, build container init cmd. exec: %s, cmd: %s",initCmd,"init")
	cmd := exec.Command(initCmd, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 如果后台启动，重定向日志输出到本地文件
		logFile,err := container.CreateLogFile(containerName)
		if err != nil {
			return nil,nil,err
		}
		cmd.Stdout = logFile
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), envSlice...)
	workDir,err := overlay.NewWorkSpace(volume, imageName, containerName)
	if err != nil {
		return nil,nil,err
	}

	// 这个决定了init程序的执行目录，方便以后的mount操作
	// 执行用户程序的时候可以设置该程序在哪个目录下执行.
	cmd.Dir = workDir
	return cmd, writePipe,nil
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
