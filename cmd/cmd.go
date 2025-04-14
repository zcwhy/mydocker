package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "mydocker",
	// 命令的简短描述
	Short: "mydocker is a simple container runtime implementation.",
}

func InitCmd() {
	rootCmd.AddCommand(NewRunCmd())
	rootCmd.AddCommand(NewInitCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
