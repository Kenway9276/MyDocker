package layer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)



func NewWorkSpace(rootDir string) string{
	createReadOnlyLayer(rootDir)
	createReadWriteLayer(rootDir)
	return createMntPoint(rootDir)
}

func DeleteWorkSpace(rootDir string) {
	fmt.Println("delete work space")
	deleteMountPoint(rootDir)
	deleteWriteLayer(rootDir)
}

func createReadOnlyLayer(rootDir string) {
	if _, err := os.Stat(path.Join(rootDir, "buzybox")); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("cannot find buzybox")
		}
	}
}

func createReadWriteLayer(rootDir string) {
	err := os.Mkdir(path.Join(rootDir, "writeLayer"), 0777)
	if err != nil {
		log.Fatal("cannot create write layer")
	}
}

func createMntPoint(rootDir string) string {
	mntDir := path.Join(rootDir, "mnt")
	err := os.Mkdir(mntDir, 0777)
	if err != nil {
		log.Fatal("cannot create mnt")
	}
	dirs := "dirs=" + path.Join(rootDir, "writeLayer") + ":" +
		path.Join(rootDir, "buzybox")
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return mntDir
}

func deleteMountPoint(rootDir string) {
	mntDir := path.Join(rootDir, "mnt")

	cmd := exec.Command("umount", mntDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(mntDir)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteWriteLayer(rootDir string) {
	writeDir := path.Join(rootDir, "writeLayer")

	err := os.RemoveAll(writeDir)
	if err != nil {
		log.Fatal(err)
	}
}