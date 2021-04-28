package volume

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	src string
	dist string
)

// mntDir: aufs mnt
func MountVolume(mntDir, volume string) error {
	ok := parseVolume(volume)
	if ok != true {
		return fmt.Errorf("volume format error")
	}

	distMount := path.Join(mntDir, dist)

	err := os.MkdirAll(distMount, 0777)
	if err != nil {
		return fmt.Errorf("cannot create guest mount point")
	}

	fmt.Printf("mount %s to %s in container\n", src, distMount)
	dirs := "dirs=" + src
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", distMount)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func UmountVolume(mntDir string) error{
	fmt.Println("umount guest")
	cmd := exec.Command("umount", path.Join(mntDir, dist))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func parseVolume(volume string) (ok bool) {
	volumeElements := strings.Split(volume, ":")
	if len(volumeElements) != 2 {
		return false
	}
	if volumeElements[0] == "" || volumeElements[1] == "" {
		return false
	}
	src = volumeElements[0]
	dist = volumeElements[1]
	return true
}