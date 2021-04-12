package command

import (
	"MyDocker/container"
	"MyDocker/resource/interface"
	"fmt"
	"github.com/urfave/cli"
	"log"
)

func Run(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("missing command")
	}
	log.Println("run command")
	var cmdArray []string
	for _, args := range ctx.Args(){
		cmdArray = append(cmdArray, args)
	}
	tty := ctx.Bool("it")
	memoryLimit := ctx.String("m")

	config := _interface.Config{
		MemoryLimit: memoryLimit,
	}
	container.RunParent(tty, cmdArray, &config)
	return nil
}

func Init(ctx *cli.Context) error {
	err := container.RunContainerInitProcess()
	return err
}


