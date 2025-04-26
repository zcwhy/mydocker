package cmd

import (
	"mydocker/container"
	"mydocker/log"
	"os"
	"os/exec"
	"syscall"

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
	firstCmd := args[0]
	path, err := exec.LookPath(firstCmd)
	if err != nil {
		log.Errorf("[initCmdRunFunc]failed to find cmd:%s", firstCmd)
		return
	}

	if err := container.SetUpMount(); err != nil {
		log.Error("[initProcess] set cotainer mount err")
		return
	}

	if err := syscall.Exec(path, args[0:], os.Environ()); err != nil {
		log.Errorf("[initProcess] exec cmd %s, err:%s", firstCmd, err)
		return
	}
}
