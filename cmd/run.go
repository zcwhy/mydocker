package cmd

import (
	"encoding/json"
	"fmt"
	"mydocker/config"
	"mydocker/container"
	"mydocker/log"
	"mydocker/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type ContainerInfo struct {
	Pid        int    `json:"pid"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	Command    string `json:"cmd"`
}

const (
	containerInfoDir  = "/home/zcy/mydocker/"
	containerInfoFile = "config.json"
)

var runOption config.RunOptions

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Create a container.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Error("Syntax Wrong with run cmd, there are no first cmd\n")
				return
			}
			runOption.CmdParams = args
			Run()
		},
	}

	runCmd.Flags().BoolVarP(&runOption.IsInteractive, "tty", "t", false, "enable tty.")
	runCmd.Flags().BoolVarP(&runOption.IsDeatch, "death", "d", false, "run container in death mod.")
	runCmd.Flags().StringVar(&runOption.Name, "name", "", "specify container name.")

	return runCmd
}

func Run() {
	createCmd, err := container.CreateContainer(&runOption)
	if err != nil {
		os.Exit(-1)
	}

	if err := createCmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	recordContainerInfo(createCmd.Process.Pid, runOption.Name, runOption.CmdParams)

	if runOption.IsInteractive {
		createCmd.Wait()
		container.DeleteWorkSpace()
	}

	os.Exit(0)
}

func recordContainerInfo(pid int, containerName string, cmd []string) {
	containerInfo := &ContainerInfo{
		Pid:     pid,
		Id:      util.GenContainerId(),
		Name:    containerName,
		Command: strings.Join(cmd, " "),
		Status:  container.STATUS_RUNNING,
	}

	if len(containerName) == 0 {
		containerInfo.Name = containerInfo.Id
	}

	infoBytes, err := json.Marshal(containerInfo)
	if err != nil {

	}

	dirPath := containerInfoDir + containerInfo.Id + "/"
	if err := os.MkdirAll(dirPath, 0677); err != nil {

	}

	filePath := dirPath + containerInfoFile
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {

	}

	if _, err := file.WriteString(string(infoBytes)); err != nil {

	}

}
