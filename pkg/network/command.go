package network

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var driver string
var subnet string

func Command() *cobra.Command {
	networkCmd := &cobra.Command{
		Use:   "network",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list container network",
		Run: func(cmd *cobra.Command, args []string) {
			ListNetwork()
		},
	}
	networkCmd.AddCommand(listCmd)

	var createCmd = &cobra.Command{
		Use:   "create [network name]",
		Short: "create a container network",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debugf("create network: %v, driver: %s, subnet: %s", args,driver,subnet)
			err := CreateNetwork(driver,subnet,args[0])
			return errors.Wrapf(err,"create network fail.")
		},
	}
	createCmd.Flags().StringVarP(&driver, "driver", "", "bridge", "network driver")
	createCmd.Flags().StringVarP(&subnet, "subnet", "", "10.0.1.0/24", "eg: 10.0.0.0/24")
	//createCmd.MarkFlagRequired("driver")
	//createCmd.MarkFlagRequired("subnet")

	networkCmd.AddCommand(createCmd)

	var deleteCmd = &cobra.Command{
		Use:   "delete [network name]",
		Short: "delete a container network",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debugf("delete network: %v", args)
			err := DeleteNetwork(args[0])
			return err
		},
	}
	networkCmd.AddCommand(deleteCmd)
	return networkCmd
}