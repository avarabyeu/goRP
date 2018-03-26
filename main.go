package main

import (
	"fmt"
	rp "github.com/avarabyeu/goRP/cli"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

var (
	version   = ""
	buildDate = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "goRP"
	app.Usage = "ReportPortal CLI Client"
	app.Version = fmt.Sprintf("%s (%s)", version, buildDate)
	app.Author = "Andrei Varabyeu"
	app.Email = "andrei.varabyeu@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "u, uuid",
			Usage:  "Access Token",
			EnvVar: "GORP_UUID",
		},
		cli.StringFlag{
			Name:   "p, project",
			Usage:  "ReportPortal Project Name",
			EnvVar: "GORP_PROJECT",
		},

		cli.StringFlag{
			Name:  "host",
			Usage: "ReportPortal Server Name",
		},
	}
	app.Commands = rp.RootCommand

	err := app.Run(os.Args)
	if nil != err {
		log.Fatal(err)
	}

}
