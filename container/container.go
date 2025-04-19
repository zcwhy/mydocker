package container

import (
	"os"
	"os/exec"
	"syscall"
)

func CreateContainer(tty bool, firstCmd string) (*exec.Cmd, error) {
	args := []string{"init", firstCmd}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
	}

	if tty {
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
