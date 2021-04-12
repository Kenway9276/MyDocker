package impl

import (
	"MyDocker/resource/interface"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
}

func (m MemorySubsystem) Name() string {
	return "memory"
}

func (m MemorySubsystem) Set(cgroupDir string, constrain _interface.Config) error {
	if cgroupPath, err := GetCgroupPath(cgroupDir, m.Name()); err == nil {
		if constrain.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(cgroupPath, "memory.limit_in_bytes"),
				[]byte(constrain.MemoryLimit), 0644); err != nil {
				return err
			}
		}
		return nil
	} else {
		return err
	}
}

func (m MemorySubsystem) Apply(cgroupDir string, pid int) error{
	if cgroupPath, err := GetCgroupPath(cgroupDir, m.Name()); err == nil {
		if err := ioutil.WriteFile(path.Join(cgroupPath, "tasks"),
			[]byte(strconv.Itoa(pid)), 0644); err != nil{
			return err
		}
		return nil
	} else {
		return err
	}
}

func (m MemorySubsystem) Remove(cgroupDir string) error {
	if cgroupPath, err := GetCgroupPath(cgroupDir, m.Name()); err == nil {
		return os.Remove(cgroupPath)
	}
	return nil
}