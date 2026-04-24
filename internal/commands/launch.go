package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/reportportal/goRP/v5/pkg/gorp"
	"github.com/reportportal/goRP/v5/pkg/openapi"
)

var errFilterNotProvided = errors.New("either IDs or filter must be provided")

var (
	launchCommand = &cli.Command{
		Name:     "launch",
		Usage:    "Operations over launches",
		Commands: []*cli.Command{listLaunchesCommand, mergeCommand},
	}

	listLaunchesCommand = &cli.Command{
		Name:  "list",
		Usage: "List launches",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "filter-name",
				Aliases: []string{"fn"},
				Usage:   "Filter Name",
				Sources: cli.EnvVars("FILTER_NAME"),
			},
			&cli.StringSliceFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Usage:   "Filter",
				Sources: cli.EnvVars("Filter"),
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
				Sources: cli.EnvVars("MERGE_LAUNCH_FILTER"),
			},
			&cli.StringFlag{
				Name:    "filter-name",
				Aliases: []string{"fn"},
				Usage:   "Filter Name",
				Sources: cli.EnvVars("FILTER_NAME"),
			},
			&cli.Int64SliceFlag{
				Name:    "ids",
				Usage:   "Launch IDS to Merge",
				Sources: cli.EnvVars("MERGE_LAUNCH_IDS"),
			},

			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "New Launch Name",
				Sources:  cli.EnvVars("MERGE_LAUNCH_NAME"),
				Required: true,
			},
			&cli.StringFlag{
				Name:    "t",
				Aliases: []string{"type"},
				Usage:   "Merge Type",
				Sources: cli.EnvVars("MERGE_TYPE"),
				Value:   "DEEP",
			},
		},
	}
)

func mergeLaunches(ctx context.Context, cmd *cli.Command) error {
	rpClient, cfg, err := buildClient(ctx, cmd)
	if err != nil {
		return err
	}

	ids, err := getMergeIDs(ctx, cmd, rpClient, cfg.Project)
	if err != nil {
		return err
	}
	rq := openapi.MergeLaunchesRQ{
		Name:      cmd.String("name"),
		MergeType: cmd.String("type"),
		Launches:  ids,
	}
	launchResource, _, err := rpClient.LaunchAsyncAPI.MergeLaunchesOldUuid(ctx, cfg.Project).
		MergeLaunchesRQ(rq).
		Execute()
	if err != nil {
		return fmt.Errorf("unable to merge launches: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.Writer, launchResource.Id)

	return nil
}

func listLaunches(ctx context.Context, cmd *cli.Command) error {
	rpClient, cfg, err := buildClient(ctx, cmd)
	if err != nil {
		return err
	}

	var launches *openapi.PageLaunchViewModel

	if filters := cmd.StringSlice("filter"); len(filters) > 0 {
		filter := strings.Join(filters, "&")
		launches, err = rpClient.GetLaunchesByFilterString(ctx, cfg.Project, filter)
	} else if filterName := cmd.String("filter-name"); filterName != "" {
		launches, err = rpClient.GetLaunchesByFilterName(ctx, cfg.Project, filterName)
	} else {
		launches, _, err = rpClient.LaunchAPI.GetProjectLaunches(ctx, cfg.Project).Execute()
	}
	if err != nil {
		return err
	}

	for _, launch := range launches.Content {
		_, _ = fmt.Fprintf(cmd.Writer, "%d #%d \"%s\"\n", launch.Id, launch.Number, launch.Name)
	}

	return nil
}

func getMergeIDs(
	ctx context.Context,
	cmd *cli.Command,
	rpClient *gorp.Client,
	project string,
) ([]int64, error) {
	if ids := cmd.Int64Slice("ids"); len(ids) > 0 {
		return ids, nil
	}

	var launches *openapi.PageLaunchViewModel
	var err error

	filter := cmd.String("filter")
	filterName := cmd.String("filter-name")
	switch {
	case filter != "":
		launches, err = rpClient.GetLaunchesByFilterString(ctx, project, filter)
	case filterName != "":
		launches, err = rpClient.GetLaunchesByFilterName(ctx, project, filterName)
	default:
		return nil, errFilterNotProvided
	}
	if err != nil {
		return nil, fmt.Errorf("unable to find launches by filter: %w", err)
	}

	ids := make([]int64, len(launches.Content))
	for i, l := range launches.Content {
		ids[i] = l.Id
	}

	return ids, nil
}
