package image

import (
	"errors"
	"github.com/huiwq1990/mydocker/pkg/types"
	"path"
)

var factory = make(map[string]string)

func init() {
	factory["alpine"] = "alpine-minirootfs-3.11.3-x86_64.tar.gz"
}

func GetImageTar(name string) (string, error) {
	if val, ok := factory[name]; ok {
		return path.Join(types.ImageRepository, val), nil
	}
	return "", errors.New("image not exist " + name)
}
