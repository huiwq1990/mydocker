// Copyright (c) 2017 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	createCmd.Flags().StringVarP(&driver, "driver", "", "bridge", "network driver(required)")
	createCmd.Flags().StringVarP(&subnet, "subnet", "", "10.0.1.0/24", "eg: 10.0.0.0/24(required)")
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


//var NetworkCommand = cli.Command{
//	Name:  "network",
//	Usage: "container network commands",
//	Subcommands: []cli.Command {
//		{
//			Name: "create",
//			Usage: "create a container network",
//			Flags: []cli.Flag{
//				cli.StringFlag{
//					Name:  "driver",
//					Usage: "network driver",
//				},
//				cli.StringFlag{
//					Name:  "subnet",
//					Usage: "subnet cidr",
//				},
//			},
//			Action:func(context *cli.Context) error {
//				if len(context.Args()) < 1 {
//					return fmt.Errorf("Missing network name")
//				}
//				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
//				if err != nil {
//					return fmt.Errorf("create network error: %+v", err)
//				}
//				return nil
//			},
//		},
//		{
//			Name: "list",
//			Usage: "list container network",
//			Action:func(context *cli.Context) error {
//				network.ListNetwork()
//				return nil
//			},
//		},
//		{
//			Name: "remove",
//			Usage: "remove container network",
//			Action:func(context *cli.Context) error {
//				if len(context.Args()) < 1 {
//					return fmt.Errorf("Missing network name")
//				}
//				err := network.DeleteNetwork(context.Args()[0])
//				if err != nil {
//					return fmt.Errorf("remove network error: %+v", err)
//				}
//				return nil
//			},
//		},
//	},
//}
