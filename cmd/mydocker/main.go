package main

import (
	"fmt"
	"github.com/huiwq1990/mydocker/pkg/container"
	"github.com/huiwq1990/mydocker/pkg/docs"
	"github.com/huiwq1990/mydocker/pkg/network"
	"github.com/huiwq1990/mydocker/pkg/run"
	"github.com/huiwq1990/mydocker/pkg/version"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	v := viper.New()
	var command = &cobra.Command{
	}

	command.AddCommand(run.Command())
	command.AddCommand(container.StopCommand())
	command.AddCommand(container.InitCommand())
	command.AddCommand(container.RemoveCommand())
	command.AddCommand(container.ExecCommand())
	command.AddCommand(container.LogsCommand())
	command.AddCommand(container.Command())
	command.AddCommand(network.Command())
	command.AddCommand(version.Command())
	command.AddCommand(docs.Command(v))

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
