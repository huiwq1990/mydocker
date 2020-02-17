// +build linux

package container

import (
	"fmt"
	"github.com/docker/docker/pkg/mount"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// init command的具体工作
func RunContainerInitProcess() error {
	log.Debugf("exec init cmd.")
	cmdArray,err := readUserCommand()
	if err != nil {
		return err
	}
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	err = setUpMount()
	if err != nil {
		return err
	}

	// 执行容器的启动命令
	log.Debugf("init command, start exec container command: %v",cmdArray)
	containerCmd, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	//syscall.exec会执行参数指定的命令,但是并不创建新的进程,只在当前进程空间内执行,即替换当前进程的执行内容,会重用同一个进程号PID.
	log.Debugf("init command, exec container action: %s",containerCmd)
	if err := syscall.Exec(containerCmd, cmdArray[0:], os.Environ()); err != nil {
		return err
	}
	return nil
}

// 启动参数是docker run传递过来的
func readUserCommand() ([]string,error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init action, read pipe error: %v", err)
		return nil, err
	}
	msgStr := string(msg)
	log.Debugf("init action, read from pipe, receive data: %s",msgStr)
	return strings.Split(msgStr, " "),nil
}

func setUpMount() error{
	pwd, err := os.Getwd()
	log.Debugf("init action, current location: %s",pwd)
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return err
	}

	err = chroot(pwd)
	//err = pivotRoot(pwd)
	if err != nil {
		return err
	}
	return nil

	//defaultMountFlags，实现用户进程成为1号进程
	// 如果不增加这个flag，业务进程id可能不是1
	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("mount proc err. %s",err.Error())
		return err
	}
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		log.Error("mount tmpfs err, %s",err.Error())
		return err
	}
	return nil
}

// 自己实现的，目前有问题
func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	log.Debugf("mount self")
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	log.Debugf("unshare current")
	// 参考docker
	if err := unix.Unshare(unix.CLONE_NEWNS); err != nil {
		return fmt.Errorf("Error creating mount namespace before pivot: %v", err)
	}
	if err := unix.PivotRoot(root,pivotDir); err != nil {
		return errors.Wrapf(err,"pivor root fail. %s",root)
	}

	return nil

	ucmd := exec.Command("unshare","-m")
	err := ucmd.Start()
	if err != nil {
		return err
	}

	log.Debugf("pivot root, newroot: %s, oldroot: %s",root,pivotDir)
	pivotCmd := exec.Command("pivot_root","./",".pivot_root")
	err = pivotCmd.Start()
	if err != nil {
		return err
	}
	return nil
	//os.Exit(0)


	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		log.Error(err)
		return errors.Wrapf(err,"pivot root error. root: %s, pivot: %s",root,pivotDir)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}

// 使用docker的代码（pkg/chrootarchive/chroot_linux.go），自己的没有调试成功
func chroot(path string) (err error) {
	// if the engine is running in a user namespace we need to use actual chroot
	//if rsystem.RunningInUserNS() {
	//	return realChroot(path)
	//}
	if err := unix.Unshare(unix.CLONE_NEWNS); err != nil {
		return fmt.Errorf("Error creating mount namespace before pivot: %v", err)
	}

	// make everything in new ns private
	if err := mount.MakeRPrivate("/"); err != nil {
		return err
	}

	if mounted, _ := mount.Mounted(path); !mounted {
		if err := mount.Mount(path, path, "bind", "rbind,rw"); err != nil {
			return realChroot(path)
		}
	}

	// setup oldRoot for pivot_root
	pivotDir, err := ioutil.TempDir(path, ".pivot_root")
	if err != nil {
		return fmt.Errorf("Error setting up pivot dir: %v", err)
	}

	var mounted bool
	defer func() {
		if mounted {
			// make sure pivotDir is not mounted before we try to remove it
			if errCleanup := unix.Unmount(pivotDir, unix.MNT_DETACH); errCleanup != nil {
				if err == nil {
					err = errCleanup
				}
				return
			}
		}

		errCleanup := os.Remove(pivotDir)
		// pivotDir doesn't exist if pivot_root failed and chroot+chdir was successful
		// because we already cleaned it up on failed pivot_root
		if errCleanup != nil && !os.IsNotExist(errCleanup) {
			errCleanup = fmt.Errorf("Error cleaning up after pivot: %v", errCleanup)
			if err == nil {
				err = errCleanup
			}
		}
	}()

	if err := unix.PivotRoot(path, pivotDir); err != nil {
		// If pivot fails, fall back to the normal chroot after cleaning up temp dir
		if err := os.Remove(pivotDir); err != nil {
			return fmt.Errorf("Error cleaning up after failed pivot: %v", err)
		}
		return realChroot(path)
	}
	mounted = true

	// This is the new path for where the old root (prior to the pivot) has been moved to
	// This dir contains the rootfs of the caller, which we need to remove so it is not visible during extraction
	pivotDir = filepath.Join("/", filepath.Base(pivotDir))

	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("Error changing to new root: %v", err)
	}

	// Make the pivotDir (where the old root lives) private so it can be unmounted without propagating to the host
	if err := unix.Mount("", pivotDir, "", unix.MS_PRIVATE|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("Error making old root private after pivot: %v", err)
	}

	// Now unmount the old root so it's no longer visible from the new root
	if err := unix.Unmount(pivotDir, unix.MNT_DETACH); err != nil {
		return fmt.Errorf("Error while unmounting old root after pivot: %v", err)
	}
	mounted = false

	return nil
}

func realChroot(path string) error {
	if err := unix.Chroot(path); err != nil {
		return fmt.Errorf("Error after fallback to chroot: %v", err)
	}
	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("Error changing to new root after chroot: %v", err)
	}
	return nil
}
