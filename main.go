package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"

	rp "github.com/avarabyeu/goRP/cli"
)

var (
	version   = ""
	buildDate = ""
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

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

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("error: %v", r)
		}
	}()
	err := app.Run(os.Args)
	if err != nil {
		//nolint:gocritic
		log.Fatal(err.Error())
	}
}
