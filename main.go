package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	rp "github.com/reportportal/goRP/v5/internal/commands"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: true,
		}),
	))

	app := cli.NewApp()
	app.Name = "goRP"
	app.Usage = "ReportPortal CLI Client"
	app.Version = fmt.Sprintf("%s (%s)", version, date)
	app.Authors = []*cli.Author{{
		Name:  "Andrei Varabyeu",
		Email: "andrei.varabyeu@gmail.com",
	}}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "uuid",
			Aliases: []string{"u"},
			Usage:   "Access Token",
			EnvVars: []string{"GORP_UUID"},
		},
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "ReportPortal Project Name",
			EnvVars: []string{"GORP_PROJECT"},
		},

		&cli.StringFlag{
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
	if err := app.Run(os.Args); err != nil {
		//nolint:gocritic
		log.Fatal(err.Error())
	}
}
