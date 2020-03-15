package network

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"os"
	"text/tabwriter"
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

func CreateNetwork(driver, subnet, name string) error {
	driverImpl,exist := drivers[driver]
	if !exist {
		return errors.New("driver not exit.")
	}

	if ifaceExist(name) {
		return errors.New("iface already exist.")
	}

	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	ip, err := ipAllocator.Allocate(cidr)
	logrus.Debugf("create network subnet: %s, alloc ip: %s. %v",subnet,ip,err)
	if err != nil {
		return err
	}
	cidr.IP = ip

	nw, err := driverImpl.Create(cidr.String(), name)
	if err != nil {
		return err
	}

	return nw.dump(defaultNetworkPath)
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			nw.Name,
			nw.IpRange.String(),
			nw.Driver,
		)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return
	}
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}

	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("Error Remove Network gateway ip: %s", err)
	}

	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("Error Remove Network DriverError: %s", err)
	}

	return nw.remove(defaultNetworkPath)
}

func Exist(name string) bool {
	_, ok := networks[name]
	return ok
}