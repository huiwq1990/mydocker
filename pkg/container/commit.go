package container

import (
	"github.com/huiwq1990/mydocker/pkg/types"
	"os/exec"
	"path"
)

func CommitContainer(containerName, imageName string) error{

	c,err := GetContainer(containerName)
	if err != nil{
		return err
	}
	imageTar := path.Join(types.ImageRepository,imageName + ".tar")

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", c.MountUrl, ".").CombinedOutput(); err != nil {
		return err
	}
	return nil
}
