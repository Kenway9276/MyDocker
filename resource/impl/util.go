package impl

import (
	"bufio"
	"os"
	"path"
	"strings"
)

func GetCgroupPath(dir string, name string) (cgroupPath string, err error) {
	mountPath := FindMountPath(name)
	cgroupPath = path.Join(mountPath, dir)
	_, err = os.Stat(cgroupPath)
	if err == nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(cgroupPath, 0755); err != nil {
				return "", err
			}
		} else {
			return cgroupPath, err
		}
	}
	return cgroupPath, nil
}

func FindMountPath(subsystem string) string {
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		txt := sc.Text()
		infos := strings.Split(txt, " ")
		for _, info := range infos {
			opts := strings.Split(info, ",")
			opt := opts[len(opts)-1]
			if opt == subsystem {
				return infos[4]
			}
		}
	}
	if err := sc.Err(); err != nil {
		return ""
	}
	return ""
}