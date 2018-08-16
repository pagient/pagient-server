package main

import (
	"os"
	"time"

	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/version"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:     "pagient",
		Version:  version.Version.String(),
		Usage:    "pagient server",
		Compiled: time.Now(),

		Authors: []*cli.Author{
			{
				Name:  "David Schneiderbauer",
				Email: "david.schneiderbauer@dschneiderbauer.me",
			},
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "/conf/app.ini",
				Usage:       "set config path",
				Destination: &config.Path,
			},
		},

		Before: func(c *cli.Context) error {
			return nil
		},

		Commands: []*cli.Command{
			Server(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
