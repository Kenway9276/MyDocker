package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {
	cmdArray := readCmdFromPipe()

	log.Printf("command: %s\n", cmdArray)

	if err := setUpMount(); err != nil {
		return err
	}

	path, err := exec.LookPath(cmdArray[0])
	fmt.Println("look path : ", path)
	if err != nil {
		return err
	}

	if err := syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func readCmdFromPipe() []string {
	readPipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(readPipe)
	if err != nil {
		log.Println("cannot read from pipe")
		return nil
	}
	cmdArray := string(msg)
	return strings.Split(cmdArray, " ")
}

func pivotRoot(root string) error {
	err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("mount failed")
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err = os.MkdirAll(pivotDir, 0777); err != nil {
		fmt.Println(err)
		return fmt.Errorf("mkdir failed")
	}

	fmt.Printf("root: %s, pivot dir: %s\n", root, pivotDir)

	if err = syscall.PivotRoot(root, pivotDir); err != nil {
		fmt.Println(err)
		return fmt.Errorf("pivot root failed")
	}
	if err = os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir failed")
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot dir failed")
	}
	return os.Remove(pivotDir)
}

func setUpMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd failed")
	}

	if err = pivotRoot(pwd); err != nil {
		return err
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}

	err = syscall.Mount("tmpfs", "/dev", "tmpfs",
		syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		return err
	}

	return nil
}
