package main

import (
	"MyDocker/command"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name: "init",
	Usage: "init container process",
	Action: command.Init,
}

var runCommand = cli.Command{
	Name:   "run",
	Usage:  "create a container. mydocker run -it [command]",
	Action: command.Run,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "set limitation of memory",
		},
	},
}
