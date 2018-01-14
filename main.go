package main

import (
	"fmt"
	"os"

	"github.com/CardInfoLink/log"
	"github.com/urfave/cli"

	"github.com/wonsikin/sank/parser"
)

const (
	// AppVersion the version of this app
	AppVersion = "v0.0.1"
	// AppName the name of this app
	AppName = "sank"
)

func main() {
	log.SetLevel(log.DebugLevel)
	app := cli.NewApp()
	app.Version = AppVersion
	app.Name = AppName
	app.HelpName = AppName
	app.Usage = "a utility for downloading files"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "./sankfile",
			Usage: "download urls and its name",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "./",
			Usage: "which folder will the images download",
		},
	}

	app.Action = func(c *cli.Context) error {
		err := parser.Parser(c.String("input"), c.String("output"))
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("[Error] %v", err)
	}
}
