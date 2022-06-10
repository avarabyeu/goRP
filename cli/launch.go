package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/reportportal/goRP/v5/gorp"
)

var (
	launchCommand = &cli.Command{
		Name:        "launch",
		Usage:       "Operations over launches",
		Subcommands: cli.Commands{listLaunchesCommand, mergeCommand},
	}

	listLaunchesCommand = &cli.Command{
		Name:  "list",
		Usage: "List launches",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "filter-name",
				Aliases: []string{"fn"},
				Usage:   "Filter Name",
				EnvVars: []string{"FILTER_NAME"},
			},
			&cli.StringSliceFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Usage:   "Filter",
				EnvVars: []string{"Filter"},
			},
		},
		Action: listLaunches,
	}

	mergeCommand = &cli.Command{
		Name:   "merge",
		Usage:  "Merge Launches",
		Action: mergeLaunches,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Usage:   "Launches Filter",
				EnvVars: []string{"MERGE_LAUNCH_FILTER"},
			},
			&cli.StringFlag{
				Name:    "filter-name",
				Aliases: []string{"fn"},
				Usage:   "Filter Name",
				EnvVars: []string{"FILTER_NAME"},
			},
			&cli.IntSliceFlag{
				Name:    "ids",
				Usage:   "Launch IDS to Merge",
				EnvVars: []string{"MERGE_LAUNCH_IDS"},
			},

			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "New Launch Name",
				EnvVars:  []string{"MERGE_LAUNCH_NAME"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "t",
				Aliases: []string{"type"},
				Usage:   "Merge Type",
				EnvVars: []string{"MERGE_TYPE"},
				Value:   "DEEP",
			},
		},
	}
)

func mergeLaunches(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}

	ids, err := getMergeIDs(c, rpClient)
	if err != nil {
		return err
	}
	rq := &gorp.MergeLaunchesRQ{
		Name:      c.String("name"),
		MergeType: gorp.MergeType(c.String("type")),
		Launches:  ids,
	}
	launchResource, err := rpClient.MergeLaunches(rq)
	if err != nil {
		return fmt.Errorf("unable to merge launches: %w", err)
	}

	//nolint:forbidigo //expected output
	fmt.Println(launchResource.ID)

	return nil
}

func listLaunches(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}

	var launches *gorp.LaunchPage

	if filters := c.StringSlice("filter"); len(filters) > 0 {
		filter := strings.Join(filters, "&")
		launches, err = rpClient.GetLaunchesByFilterString(filter)
	} else if filterName := c.String("filter-name"); filterName != "" {
		launches, err = rpClient.GetLaunchesByFilterName(filterName)
	} else {
		launches, err = rpClient.GetLaunches()
	}
	if err != nil {
		return err
	}

	//nolint:forbidigo //expected output
	for _, launch := range launches.Content {
		fmt.Printf("%d #%d \"%s\"\n", launch.ID, launch.Number, launch.Name)
	}

	return nil
}

func getMergeIDs(c *cli.Context, rpClient *gorp.Client) ([]int, error) {
	if ids := c.IntSlice("ids"); len(ids) > 0 {
		return ids, nil
	}

	var launches *gorp.LaunchPage
	var err error

	filter := c.String("filter")
	filterName := c.String("filter-name")
	switch {
	case filter != "":
		launches, err = rpClient.GetLaunchesByFilterString(filter)
	case filterName != "":
		launches, err = rpClient.GetLaunchesByFilterName(filterName)
	default:
		return nil, errors.New("no either IDs or filter provided")
	}
	if err != nil {
		return nil, fmt.Errorf("unable to find launches by filter: %w", err)
	}

	ids := make([]int, len(launches.Content))
	for i, l := range launches.Content {
		ids[i] = l.ID
	}

	return ids, nil
}
