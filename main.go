package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avarabyeu/goRP/gorp"
	"github.com/manifoldco/promptui"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/user"
	"path/filepath"
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

type config struct {
	UUID    string
	Project string
	Host    string
}

var (
	rootCommands = []cli.Command{
		launchCommand,
		initCommand,
	}

	launchCommand = cli.Command{
		Name:  "launch",
		Usage: "Operations over launches",
		Subcommands: cli.Commands{
			listLaunchesCommand,
		},
	}

	initCommand = cli.Command{
		Name:   "init",
		Usage:  "Initializes configuration cache",
		Action: initConfiguration,
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

func initConfiguration(c *cli.Context) error {

	if configFilePresent() {
		prompt := promptui.Select{
			Label: "GoRP is already configured. Replace existing configuration?",
			Items: []string{"No", "Yes"},
		}
		num, _, err := prompt.Run()
		if err != nil {
			return err
		}
		//do not replace. go away
		if 0 == num {
			return nil
		}
	}
	f, err := os.OpenFile(getConfigFile(), os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return cli.NewExitError(fmt.Sprintf("Cannot open config file, %s", err), 1)
	}
	defer f.Close()

	prompt := promptui.Prompt{
		Label: "ReportPortal hostname",
	}
	host, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "UUID",
	}
	uuid, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "Default Project",
	}
	project, err := prompt.Run()
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(&config{
		Project: project,
		Host:    host,
		UUID:    uuid,
	})
	if nil != err {
		return cli.NewExitError(fmt.Sprintf("Cannot read config file. %s", err), 1)
	}

	fmt.Println("Configuration has been successfully saved!")
	return nil
}

func configFilePresent() bool {
	_, err := os.Stat(getConfigFile())
	return !os.IsNotExist(err)
}

func getConfigFile() string {
	return filepath.Join(getHomeDir(), ".gorp")
}
func getHomeDir() string {
	if h := os.Getenv("HOME"); "" != h {
		return h
	}
	curUser, err := user.Current()
	if err != nil {
		// well, sheesh
		return "."
	}

	return curUser.HomeDir
}

func getConfig(c *cli.Context) (*config, error) {
	cfg := &config{}
	if configFilePresent() {
		f, err := os.Open(getConfigFile())
		if nil != err {
			return nil, err
		}
		err = json.NewDecoder(f).Decode(cfg)
		if nil != err {
			return nil, err
		}
	}
	if v := c.GlobalString("uuid"); "" != v {
		cfg.UUID = v
	}
	if v := c.GlobalString("project"); "" != v {
		cfg.Project = v
	}
	if v := c.GlobalString("host"); "" != v {
		cfg.Host = v
	}

	if err := validateConfig(cfg); nil != err {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *config) error {
	if "" == cfg.UUID {
		return errors.New("uuid is not set")
	}

	if "" == cfg.Project {
		return errors.New("project is not set")
	}

	if "" == cfg.Host {
		return errors.New("host is not set")
	}
	return nil
}

func buildClient(ctx *cli.Context) (*gorp.Client, error) {
	cfg, err := getConfig(ctx)
	if nil != err {
		return nil, err
	}
	return gorp.NewClient(cfg.Host, cfg.Project, cfg.UUID), nil

}
