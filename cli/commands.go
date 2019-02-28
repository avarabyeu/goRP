package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/avarabyeu/goRP/gorp"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

type config struct {
	UUID    string
	Project string
	Host    string
}

var (
	// RootCommand is CLI entry point
	RootCommand = []cli.Command{
		launchCommand,
		initCommand,
		mergeCommand,
	}

	initCommand = cli.Command{
		Name:   "init",
		Usage:  "Initializes configuration cache",
		Action: initConfiguration,
	}
)

func initConfiguration(c *cli.Context) error {

	if configFilePresent() {
		prompt := promptui.Prompt{
			Label: "GoRP is already configured. Replace existing configuration?",
		}
		answer, err := prompt.Run()
		if err != nil {
			return err
		}
		//do not replace. go away
		if !answerYes(answer) {
			return nil
		}
	}
	f, err := os.OpenFile(getConfigFile(), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Cannot open config file, %s", err), 1)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Error(closeErr)
		}
	}()

	prompt := promptui.Prompt{
		Label: "Enter ReportPortal hostname",
		Validate: func(host string) error {
			_, parseErr := url.Parse(host)
			return parseErr
		},
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
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Cannot read config file. %s", err), 1)
	}

	fmt.Println("Configuration has been successfully saved!")
	return nil
}

func getConfig(c *cli.Context) (*config, error) {
	cfg := &config{}
	if configFilePresent() {
		f, err := os.Open(getConfigFile())
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(f).Decode(cfg)
		if err != nil {
			return nil, err
		}
	}
	if v := c.GlobalString("uuid"); v != "" {
		cfg.UUID = v
	}
	if v := c.GlobalString("project"); v != "" {
		cfg.Project = v
	}
	if v := c.GlobalString("host"); v != "" {
		cfg.Host = v
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func buildClient(ctx *cli.Context) (*gorp.Client, error) {
	cfg, err := getConfig(ctx)
	if err != nil {
		return nil, err
	}
	return gorp.NewClient(cfg.Host, cfg.Project, cfg.UUID), nil
}
