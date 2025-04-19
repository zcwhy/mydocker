package cmd

import (
	"fmt"
	"mydocker/container"
	"os"

	"github.com/spf13/cobra"
)

var tty bool
var command string

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Create a container.",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}

	runCmd.Flags().BoolVarP(&tty, "tty", "t", false, "enable tty.")
	runCmd.Flags().StringVarP(&command, "interactive", "i", "/bin/bash", "interactive mod.")

	return runCmd
}

func Run() {
	createCmd, err := container.CreateContainer(tty, command)
	if err != nil {
		os.Exit(-1)
	}

	if err := createCmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	createCmd.Wait()

	// fmt.Println("1111111111111111")
	container.DeleteWorkSpace()

	os.Exit(-1)
}
