package container

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "create a container network",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ListContainers()
			return nil
		},
	}
	return cmd
}

func LogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			return LogContainer(containerName)
		},
	}
	return cmd
}


func StopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			stopContainer(containerName)
			return nil
		},
	}
	return cmd
}


func RemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "remove unused containers",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			removeContainer(containerName)
			return nil
		},
	}
	return cmd
}



func CommitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "commit a container into image",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			return CommitContainer(containerName, args[1])
		},
	}
	return cmd
}



func ExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "exec a command into container",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 命令回调时会执行
			if os.Getenv(ENV_EXEC_PID) != "" {
				return errors.Errorf("pid callback pid %s", os.Getgid())
			}

			ExecContainer(args[0], args[1:])
			return nil
		},
	}
	return cmd
}
