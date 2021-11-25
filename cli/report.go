package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	reportCommand = &cli.Command{
		Name:        "report",
		Usage:       "Reports input to report portal",
		Subcommands: cli.Commands{listLaunchesCommand, mergeCommand},
	}

	reportTest2JsonCommand = &cli.Command{
		Name:   "test2json",
		Usage:  "Input format: test2json",
		Action: reportTest2json,
	}
)

func reportTest2json(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}

	fmt.Println(rpClient)
	return nil
}
