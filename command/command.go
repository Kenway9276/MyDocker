package command

import (
	"MyDocker/container"
	"MyDocker/layer"
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
	detach := ctx.Bool("d")
	memoryLimit := ctx.String("m")

	if tty && detach {
		return fmt.Errorf("tty or detach")
	}

	config := _interface.Config{
		MemoryLimit: memoryLimit,
	}

	volume := ctx.String("v")
	err := container.RunParent(tty, cmdArray, &config, volume)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Init(_ *cli.Context) error {
	log.Println("run init")
	err := container.RunContainerInitProcess()
	return err
}

func Commit(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("missing command")
	}
	imageName := ctx.Args().Get(0)
	layer.CommitContainer(imageName)
	return nil
}
