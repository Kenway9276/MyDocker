package _interface

type Subsystem interface {
	Name() string
	Set(cgroupDir string, constrain Config) error
	Apply(cgroupDir string, pid int) error
	Remove(cgroupDir string) error
}

type Config struct {
	MemoryLimit string
	CpuShare string
	CpuSet string
}