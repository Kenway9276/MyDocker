package container

import (
	"MyDocker/layer"
	"MyDocker/resource"
	_interface "MyDocker/resource/interface"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func RunParent(tty bool, cmdArray []string, config *_interface.Config) {

	parent := newParentPrecess(tty, cmdArray)

	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}

	subsystemManager := resource.NewSubsystemManager("mydocker-cgroup")
	subsystemManager.Set(*config)
	subsystemManager.Apply(parent.Process.Pid)

	parent.Wait()
}

func newParentPrecess(tty bool, cmdArray []string) *exec.Cmd {
	resCmd := exec.Command(`/proc/self/exe`, "init")

	resCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID,
	}
	if tty {
		resCmd.Stdin = os.Stdin
		resCmd.Stdout = os.Stdout
		resCmd.Stderr = os.Stderr
	}

	reader, writer, err := NewPipe()
	if err != nil {
		log.Fatal(err)
	}

	writeCmdToPipe(cmdArray, writer)
	resCmd.ExtraFiles = []*os.File{reader}
	resCmd.Dir = layer.NewWorkSpace("/mydocker")
	return resCmd
}

func writeCmdToPipe(cmdArray []string, writer *os.File) {
	defer writer.Close()
	cmd := strings.Join(cmdArray, " ")
	_, err := io.WriteString(writer, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func NewPipe() (reader *os.File, writer *os.File, err error) {
	return os.Pipe()
}