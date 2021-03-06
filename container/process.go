package container

import (
	"MyDocker/layer"
	"MyDocker/resource"
	_interface "MyDocker/resource/interface"
	"MyDocker/volume"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var workDir = "/mydocker"
var MntDir string

func RunParent(tty bool, cmdArray []string, config *_interface.Config, volumeCmd string) error{
	MntDir = layer.NewWorkSpace(workDir) // workDir/mnt
	if tty {
		defer layer.DeleteWorkSpace(workDir)
	}

	if volumeCmd != ""{
		if err := volume.MountVolume(MntDir, volumeCmd); err != nil {
			return err
		}
		if tty {
			defer volume.UmountVolume(MntDir)
		}
	}

	parent, err := newParentPrecess(tty, cmdArray, MntDir)
	if err != nil {
		return err
	}

	if err := parent.Start(); err != nil {
		log.Println(err)
		return err
	}

	subsystemManager := resource.NewSubsystemManager("mydocker-cgroup")
	subsystemManager.Set(*config)
	subsystemManager.Apply(parent.Process.Pid)

	if tty {
		parent.Wait()
	}
	return nil
}

func newParentPrecess(tty bool, cmdArray []string, workDir string) (*exec.Cmd, error) {
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
	resCmd.Dir = workDir

	return resCmd, nil
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