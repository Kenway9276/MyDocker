package resource

import (
	"MyDocker/resource/impl"
	"MyDocker/resource/interface"
)

type SubsystemManager struct {
	cgroupDir string
}

var (
	SubSystemIns = []_interface.Subsystem{
		&impl.MemorySubsystem{},
	}
)

func (m SubsystemManager) Apply(pid int) (err error) {
	for _, subSystemIns := range SubSystemIns {
		err = subSystemIns.Apply(m.cgroupDir, pid)
	}
	return
}

func (m SubsystemManager) Set(config _interface.Config) (err error) {
	for _, subsystemIns := range SubSystemIns {
		err = subsystemIns.Set(m.cgroupDir, config)
	}
	return
}

func (m SubsystemManager) Destroy() (err error) {
	for _, subsystemIns := range SubSystemIns {
		err = subsystemIns.Remove(m.cgroupDir)
	}
	return
}

func NewSubsystemManager(cgroupDir string) SubsystemManager {
	return SubsystemManager{cgroupDir: cgroupDir}
}