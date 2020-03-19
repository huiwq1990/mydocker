package nsenter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestGetCurrentThreadNSPath(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentPID := os.Getpid()
	currentTID := unix.Gettid()
	for nsType := range CloneFlagsTable {
		expectedPath := fmt.Sprintf("/proc/%d/task/%d/ns/%s", currentPID, currentTID, nsType)
		path := getCurrentThreadNSPath(nsType)
		assert.Equal(t, path, expectedPath)
	}
}


func TestNsEnterSuccessful(t *testing.T) {
	if tu.NotValid(ktu.NeedRoot()) {
		t.Skip(ktu.TestDisabledNeedRoot)
	}
	nsList := supportedNamespaces()
	sleepDuration := 60

	cloneFlags := 0
	for _, ns := range nsList {
		cloneFlags |= CloneFlagsTable[ns.Type]
	}

	sleepPID, err := startSleepBinary(sleepDuration, cloneFlags)
	assert.NoError(t, err)
	defer func() {
		if sleepPID > 1 {
			unix.Kill(sleepPID, syscall.SIGKILL)
		}
	}()

	for idx := range nsList {
		nsList[idx].Path = getNSPathFromPID(sleepPID, nsList[idx].Type)
		nsList[idx].PID = sleepPID
	}

	var sleepPIDFromNsEnter int

	testToRun := func() error {
		sleepPIDFromNsEnter, err = startSleepBinary(sleepDuration, 0)
		if err != nil {
			return err
		}

		return nil
	}

	err = NsEnter(nsList, testToRun)
	assert.Nil(t, err, "%v", err)

	defer func() {
		if sleepPIDFromNsEnter > 1 {
			unix.Kill(sleepPIDFromNsEnter, syscall.SIGKILL)
		}
	}()

	for _, ns := range nsList {
		nsPathEntered := getNSPathFromPID(sleepPIDFromNsEnter, ns.Type)

		// Here we are trying to resolve the path but it fails because
		// namespaces links don't really exist. For this reason, the
		// call to EvalSymlinks will fail when it will try to stat the
		// resolved path found. As we only care about the path, we can
		// retrieve it from the PathError structure.
		evalExpectedNSPath, err := filepath.EvalSymlinks(ns.Path)
		if err != nil {
			evalExpectedNSPath = err.(*os.PathError).Path
		}

		// Same thing here, resolving the namespace path.
		evalNSEnteredPath, err := filepath.EvalSymlinks(nsPathEntered)
		if err != nil {
			evalNSEnteredPath = err.(*os.PathError).Path
		}

		_, evalExpectedNS := filepath.Split(evalExpectedNSPath)
		_, evalNSEntered := filepath.Split(evalNSEnteredPath)

		assert.Equal(t, evalExpectedNS, evalNSEntered)
	}
}
