package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"mydocker/container"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "the first process in which container create exec.",
		Run:   initCmdRunFunc,
	}

	return initCmd
}

func initCmdRunFunc(cmd *cobra.Command, args []string) {
	fmt.Println("init called: Hello mydocker")
	container.SetUpMount()

	firstCmd := args[0]
	path, err := exec.LookPath(firstCmd)
	if err != nil {

	}

	fmt.Printf("first cmd: %s\n", firstCmd)
	if err := syscall.Exec(path, args[0:], os.Environ()); err != nil {

	}
}
