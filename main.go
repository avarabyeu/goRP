package main

import (
	"fmt"
	"github.com/avarabyeu/goRP/gorp"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strings"
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
	app.Commands = rootCommands

	err := app.Run(os.Args)
	if nil != err {
		log.Fatal(err)
	}

}

func buildClient(c *cli.Context) (*gorp.Client, error) {
	uuid, err := requiredFlag("uuid", c)
	if nil != err {
		return nil, err
	}
	proj, err := requiredFlag("project", c)
	if nil != err {
		return nil, err
	}
	host, err := requiredFlag("host", c)
	if nil != err {
		return nil, err
	}

	return gorp.NewClient(host, proj, uuid), nil

}

func requiredFlag(f string, c *cli.Context) (string, error) {
	fVal := c.GlobalString(f)
	if "" == fVal {
		return "", cli.NewExitError(fmt.Sprintf("%s is not set", f), 1)
	}
	return fVal, nil
}

var (
	rootCommands = []cli.Command{
		launchesCommand,
	}

	launchesCommand = cli.Command{
		Name:  "launch",
		Usage: "Operations over launches",
		Subcommands: cli.Commands{
			listLaunchesCommand,
		},
	}

	listLaunchesCommand = cli.Command{
		Name:  "list",
		Usage: "List launches",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "fn, filter-name",
				Usage:  "Filter Name",
				EnvVar: "FILTER_NAME",
			},
			cli.StringSliceFlag{
				Name:   "f, filter",
				Usage:  "Filter",
				EnvVar: "Filter",
			},
		},
		Action: listLaunches,
	}
)

func listLaunches(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if nil != err {
		return err
	}

	var launches *gorp.LaunchPage

	if filters := c.StringSlice("filter"); nil != filters && len(filters) > 0 {
		filter := strings.Join(filters, "&")
		launches, err = rpClient.GetLaunchesByFilterString(filter)
	} else if filterName := c.String("filter-name"); "" != filterName {
		launches, err = rpClient.GetLaunchesByFilterName(filterName)
	} else {
		launches, err = rpClient.GetLaunches()
	}
	if nil != err {
		return err
	}

	for _, launch := range launches.Content {
		fmt.Printf("%s #%d \"%s\"\n", launch.ID, launch.Number, launch.Name)
	}
	return nil
}
