package container

import (
	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "container init",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := RunContainerInitProcess()
			return err
		},
	}

	return cmd
}