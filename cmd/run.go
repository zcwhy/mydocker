package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mydocker/container"
	"os"
)

var tty bool
var command string

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Create a container.",
		Run: func(cmd *cobra.Command, args []string) {
			Run(tty, command)
		},
	}

	runCmd.Flags().BoolVarP(&tty, "tty", "t", false, "enable tty.")
	runCmd.Flags().StringVarP(&command, "interactive", "i", "/bin/bash", "interactive mod.")

	return runCmd
}

func Run(tty bool, cmd string) {
	createCmd := container.CreateContainer(tty)

	if err := createCmd.Start(); err != nil {
		fmt.Println(err)
	}
	createCmd.Wait()
	os.Exit(-1)
}
