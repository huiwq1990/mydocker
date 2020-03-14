package overlay

import (
	"errors"
	"github.com/huiwq1990/mydocker/pkg/image"
	"github.com/huiwq1990/mydocker/pkg/types"
	"github.com/huiwq1990/mydocker/pkg/util"
	log "github.com/sirupsen/logrus"
	"path"

	"os"
	"os/exec"
	"strings"
)

func NewWorkSpace(volume, imageName, containerName string) (string,error) {
	err := createReadOnlyLayer(imageName)
	if err != nil {
		return "",err
	}


	workDir,err := doOver(imageName,containerName)
	if err != nil {
		return "",err
	}


	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			if err := MountVolume(volumeURLs[0],volumeURLs[1], containerName); err != nil {
				return "",err
			}
		} else {
			return "",errors.New("Volume parameter input is not correct.")
		}
	}

	//TODO 需要结合overlay的upper worker层看看怎么实现

	return workDir,nil
}

func doOver(imageName,containerName string) (string,error){
	containerRootDir := path.Join(types.WriteLayerUrl, containerName)
	log.Infof("create write layer. %s",containerRootDir)
	if err := os.MkdirAll(containerRootDir, 0777); err != nil {
		log.Errorf("Mkdir write layer dir %s error. %v", containerRootDir, err)
		return "", err
	}

	upperDir := path.Join(containerRootDir,"upper")
	log.Debugf("create write layer. %s",upperDir)
	if err := os.MkdirAll(upperDir, 0777); err != nil {
		log.Errorf("Mkdir write layer dir %s error. %v", upperDir, err)
		return "", err
	}

	workDir := path.Join(containerRootDir,"worker")
	log.Debugf("create write layer. %s",workDir)
	if err := os.MkdirAll(workDir, 0777); err != nil {
		log.Errorf("Mkdir write layer dir %s error. %v", workDir, err)
		return "", err
	}

	mergeDir := path.Join(containerRootDir, "merge")
	log.Infof("create write layer. %s",mergeDir)
	if err := os.MkdirAll(mergeDir, 0777); err != nil {
		log.Errorf("Mkdir write layer dir %s error. %v", mergeDir, err)
		return "", err
	}

	//mount -t overlay overlay -o lowerdir=/root/busybox,upperdir=upper,workdir=worker merge
	oDir := "lowerdir=/root/"+imageName+",upperdir="+ upperDir +",workdir="+workDir
	log.Debugf("mount -t overlay overlay -o %s %s",oDir,mergeDir)
	out, err := exec.Command("mount", "-t", "overlay","overlay", "-o", oDir, mergeDir).CombinedOutput()
	if err != nil {
		//'Special device overlay doesn't exist' 可能是目录不存在
		log.Errorf("mount overlay failed, output:%s, err:%v", string(out), err)
		return "", err
	}
	return mergeDir, err
}

func MountVolume(from string, target string, containerName string) error {
	targetMountDir := path.Join(types.WriteLayerUrl, containerName,"merge",target)
	if err := os.MkdirAll(targetMountDir, 0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v", targetMountDir, err)
	}
	log.Debugf("mount --bind %s %s",from,targetMountDir)
	_, err := exec.Command("mount", "--bind", from, targetMountDir).CombinedOutput()
	if err != nil {
		log.Errorf("Mount volume failed. %v", err)
		return err
	}

	out, err := exec.Command("mount","-o","remount,rw,bind",targetMountDir).CombinedOutput()
	if err != nil {
		log.Errorf("Mount volume failed.%v %v",string(out), err)
		return err
	}

	return nil
}

func DeleteWorkSpace(volume, containerName string) error{
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			if _, err := exec.Command("umount", path.Join(types.WriteLayerUrl,containerName,"upper",volumeURLs[0])).CombinedOutput(); err != nil {
				return err
			}
		}else{
			return errors.New("mount config error" + volume)
		}
	}
	return DeleteMountPoint(containerName)
}

func DeleteMountPoint(containerId string) error {

	mountDir := path.Join(types.WriteLayerUrl,containerId, "merge")
	_, err := exec.Command("umount", mountDir).CombinedOutput()
	if err != nil {
		log.Errorf("Unmount %s error %v", mountDir, err)
		return err
	}
	if err := os.RemoveAll(path.Join(types.WriteLayerUrl,containerId)); err != nil {
		return err
	}
	return nil
}

func createReadOnlyLayer(imageName string) error {
	untarFolderUrl := path.Join(types.ImageRepository, imageName)
	exist, err := util.PathExists(untarFolderUrl)
	if err != nil {
		log.Errorf("Fail to judge whether dir %s exists. %v", imageName, err)
		return err
	}
	if !exist {
		tarUrl, err := image.GetImageTar(imageName)
		if err != nil{
			return err
		}
		log.Debugf("crate read only layer. image:%s, untarurl:%s",untarFolderUrl,tarUrl)
		if err := os.MkdirAll(untarFolderUrl, 0622); err != nil {
			log.Errorf("Mkdir %s error %v", untarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", tarUrl, "-C", untarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", tarUrl, err)
			return err
		}
	}
	return nil
}
