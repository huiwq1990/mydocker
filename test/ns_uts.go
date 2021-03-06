// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main()  {
	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
	}
	cmd.Stdin  = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Run error:%v\n", err)
		log.Fatal(err)
	}
}