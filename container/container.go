package container

import (
	"fmt"
	"mydocker/config"
	"os"
	"os/exec"
	"syscall"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	Command    string `json:"cmd"`
}

const (
	STATUS_RUNNING string = "running"
	STATUS_STOP    string = "stop"
	STATUS_EXIT    string = "exited"
)

const ()

func CreateContainer(cmdOptions *config.RunOptions) (*exec.Cmd, error) {
	initParams := append([]string{"init"}, cmdOptions.CmdParams...)
	// 调用自身可执行文件，执行init命令
	cmd := exec.Command("/proc/self/exe", initParams...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
	}

	if cmdOptions.IsDeatch && cmdOptions.IsInteractive {
		return nil, fmt.Errorf("id and d paramter can not both provided")
	}

	if cmdOptions.IsInteractive {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	mntURL := "/home/zcy/mnt"
	rootURL := "/home/zcy/"
	err := CreateWorkSpace(rootURL, mntURL)
	if err != nil {
		return nil, err
	}
	cmd.Dir = mntURL

	return cmd, nil
}

func CreateWorkSpace(baseUrl string, mntUrl string) error {
	lowerdir, err := createReadOnlyLayer(baseUrl)
	if err != nil {
		return err
	}

	upperdir, err := createWriteLayer(baseUrl)
	if err != nil {
		return err
	}

	return createMountPoint(lowerdir, upperdir, mntUrl)
}
