package main

import (
	"os"
	"syscall"
	"testing"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"

)

func ss()  {
	//graphdriver
}

func cdMountFrom(dir, device, target, mType, label string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(dir)
	defer os.Chdir(wd)

	return syscall.Mount(device, target, mType, 0, label)
}

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestOverlaySetup and TestOverlayTeardown
func TestOverlaySetup(t *testing.T) {
	graphtest.GetDriver(t, driverName)
}
