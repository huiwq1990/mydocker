package container

import (
	"fmt"
	"github.com/huiwq1990/mydocker/pkg/types"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

func ListContainers() error {

	files, err := ioutil.ReadDir(types.WriteLayerUrl)
	if err != nil {
		return err
	}

	var containers []*types.ContainerInfo
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		tmpContainer, err := GetContainer(file.Name())
		if err != nil {
			return err
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

