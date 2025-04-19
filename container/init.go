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
func SetUpMount() error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
	}

	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		log.Errorf("")
		return err
	}

	err = pivotRoot(cwd)
	if err != nil {
		log.Errorf("[SetUpMount] pivot_root err: %s", err)
		return err
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		log.Errorf("[SetUpMount] mount proc to /proc err: %s", err)
		return err
	}

	return nil
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.MkdirAll(pivotDir, 0777); err != nil {
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

	exist, err := util.PathExists(bustBoxTarUrl)
	if err != nil || !exist {
		log.Errorf("[createReadOnlyLayer] find path %s err %v", bustBoxTarUrl, err)
		return "", err
	}

	if err := os.MkdirAll(busyBoxUrl, 0777); err != nil {
		log.Errorf("[createReadOnlyLayer] mkdir %s err %v", busyBoxUrl, err)
		return "", err
	}

	if _, err := exec.Command("tar", "-xvf", bustBoxTarUrl, "-C", busyBoxUrl).CombinedOutput(); err != nil {
		log.Errorf("[createReadOnlyLayer] unTar dir %s err %v", bustBoxTarUrl, err)
		return "", err
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

	workDir := "/home/zcy/work"
	if err := os.MkdirAll(workDir, 0777); err != nil {
		log.Errorf("[createMountPoint] create workDir: %s, err: %v", workDir, err)
		return err
	}
	args := "lowerdir=" + lowerdir + ",upperdir=" + upperdir + ",workdir=" + workDir
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", args, mntUrl)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("[createMountPoint] mount overlayfs arg:%s, mnt point:%s err: %v", args, mntUrl, err)
		return err
	}

	return nil
}

// 容器结束后清理工作空间
func DeleteWorkSpace() {
	rootUrl := "/home/zcy/"
	mntUrl := "/home/zcy/mnt/"
	DeleteMountPoint(mntUrl)

	DeleteWriteLayer(rootUrl)
}

func DeleteMountPoint(mntUrl string) {
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {

	}

	if err := os.RemoveAll(mntUrl); err != nil {

	}
}

func DeleteWriteLayer(rootUrl string) {
	writeUrl := rootUrl + "writeLayer/"
	if err := os.RemoveAll(writeUrl); err != nil {

	}

	workUrl := rootUrl + "work/"
	if err := os.RemoveAll(workUrl); err != nil {

	}
}
