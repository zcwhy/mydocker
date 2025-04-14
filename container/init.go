package container

import (
	"fmt"
	"mydocker/log"
	"mydocker/util"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// 设置文件系统
func SetUpMount() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Get current location error %v", err)
	}
	fmt.Printf("Current location is %s\n", cwd)

	entries, err := os.ReadDir(cwd)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return
	}

	fmt.Println("目录内容：")
	for _, entry := range entries {
		// 打印所有文件和文件夹，包括隐藏文件（以 . 开头）
		fmt.Println(entry.Name())
	}
	// err = pivotRoot(cwd)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // mount proc
	// syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "")
	// syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
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

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v, pivot_root: %s", err, pivotDir)
	}

	return os.Remove(pivotDir)
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

func createReadOnlyLayer(baseUrl string) (string, error) {
	busyBoxUrl := baseUrl + "busybox/"
	bustBoxTarUrl := baseUrl + "busybox.tar"

	exist, err := util.PathExists(busyBoxUrl)
	if err != nil {
		log.Errorf("[createReadOnlyLayer] find path %s err %v", bustBoxTarUrl, err)
		return "", err
	}

	if !exist {
		if err := os.Mkdir(busyBoxUrl, 0777); err != nil {
			log.Errorf("[createReadOnlyLayer] mkdir %s err %v", busyBoxUrl, err)
			return "", err
		}

		if _, err := exec.Command("tar", "-xvf", bustBoxTarUrl, "-C", busyBoxUrl).CombinedOutput(); err != nil {
			log.Errorf("[createReadOnlyLayer] unTar dir %s err %v", bustBoxTarUrl, err)
			return "", err
		}
	}
	return busyBoxUrl, nil
}

func createWriteLayer(baseUrl string) (string, error) {
	writeURL := baseUrl + "writeLayer/"
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Errorf("[createWriteLayer] make writelayer %s err:%v", writeURL, err)
		return "", err
	}
	return writeURL, nil
}

// use overlay
func createMountPoint(lowerdir, upperdir, mntUrl string) error {
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Errorf("[createMountPoint] create mount dir: %s, err: %v", mntUrl, err)
		return err
	}

	workDir := "/root/work"
	if err := os.MkdirAll(workDir, 0777); err != nil {
		log.Errorf("[createMountPoint] create workDir: %s, err: %v", workDir, err)
		return err
	}
	args := "lowerdir=" + lowerdir + ",upperdir=" + upperdir + ",workdir=" + workDir
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", args, mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("[createMountPoint] mount aufs arg:%s, mnt point:%s err: %v", args, mntUrl, err)
		return err
	}

	return nil
}
