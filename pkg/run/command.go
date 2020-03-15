package run

import (
	"fmt"
	"github.com/huiwq1990/mydocker/pkg/network"
	"github.com/huiwq1990/mydocker/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "create a container",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			imageName := args[0]
			cmdArray := args[1:]

			config.ImageName = imageName
			config.CmdArray = cmdArray

			logrus.Debugf("run action, image: %s, exec args: %s",imageName, cmdArray)
			// 不能即是不前台进程，又不是后台进程
			if config.Tty && config.Detach {
				return fmt.Errorf("ti and d paramter can not both provided")
			}

			if config.Net != "" {
				if !network.Exist(config.Net){
					return fmt.Errorf("network %v not exist",config.Net)
				}
			}

			resConf := &types.ResourceConfig{
				MemoryLimit: config.Memory,
				CpuSet:      config.Cpuset,
				CpuShare:    config.Cpushare,
			}
			logrus.Debugf("run config: %s", config.String())

			return Run(config, resConf)
		},
	}

	//cmd.Flags().StringVarP(&config.Name, "name", "", "", "name(required)")
	cmd.Flags().StringVarP(&config.Net, "net", "n", "", "net")
	cmd.Flags().BoolVarP(&config.Tty, "tty", "t", true, "tty")
	cmd.Flags().BoolVarP(&config.Detach, "detach", "d", false, "detach")

	cmd.Flags().StringVarP(&config.Memory, "m", "", "", "memory limit")
	cmd.Flags().StringVarP(&config.Cpushare, "cpushare", "", "", "cpushare limit")
	cmd.Flags().StringVarP(&config.Cpuset, "cpuset", "", "", "cpuset limit")
	cmd.Flags().StringVarP(&config.Volume, "volume", "v", "", "mount volume")

	return cmd
}

var config = types.RunCommandConfig{}

