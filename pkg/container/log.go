package container

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func LogContainer(containerName string) error {

	c,err := GetContainer(containerName)
	if err != nil {
		return err
	}

	file, err := os.Open(c.LogFile)
	defer file.Close()
	if err != nil {
		return errors.Wrapf(err,"open log file fail.")
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, string(content))
	return nil
}
