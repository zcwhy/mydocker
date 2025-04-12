package container

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// 设置文件系统
func SetUpMount() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Get current location error %v", err)
	}
	fmt.Sprint("Current location is %s", cwd)

	pivotRoot(cwd)
	// mount proc
	syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err) //删除临时文件夹
	}

	return os.Remove(pivotDir)
}
