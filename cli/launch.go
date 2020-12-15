package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/avarabyeu/goRP/gorp"
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
				Name:    "fn, filter-name",
				Usage:   "Filter Name",
				EnvVars: []string{"FILTER_NAME"},
			},
			&cli.StringSliceFlag{
				Name:    "f, filter",
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
				Name:    "f, filter",
				Usage:   "Launches Filter",
				EnvVars: []string{"MERGE_LAUNCH_FILTER"},
			},
			&cli.StringFlag{
				Name:    "fn, filter-name",
				Usage:   "Filter Name",
				EnvVars: []string{"FILTER_NAME"},
			},
			&cli.StringSliceFlag{
				Name:    "ids",
				Usage:   "Launch IDS to Merge",
				EnvVars: []string{"MERGE_LAUNCH_IDS"},
			},

			&cli.StringFlag{
				Name:    "n, name",
				Usage:   "New Launch Name",
				EnvVars: []string{"MERGE_LAUNCH_NAME"},
			},
			&cli.StringFlag{
				Name:    "t, type",
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

	for _, launch := range launches.Content {
		fmt.Printf("%s #%d \"%s\"\n", launch.ID, launch.Number, launch.Name)
	}

	return nil
}

func getMergeIDs(c *cli.Context, rpClient *gorp.Client) ([]string, error) {
	if ids := c.StringSlice("ids"); len(ids) > 0 {
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
		return nil, fmt.Errorf("unable to find launches by filter: %s", err.Error())
	}

	ids := make([]string, len(launches.Content))
	for i, l := range launches.Content {
		ids[i] = l.ID
	}

	return ids, nil
}
