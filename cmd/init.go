package cmd

import (
	"fmt"

	"mydocker/container"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "the first process in which container create exec.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("init called: %s\n", args[0])
			InitContainer()
		},
	}

	return initCmd
}

func InitContainer() {
	container.SetUpMount()
}
