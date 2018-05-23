package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	cfg := config.New()

	if env := os.Getenv("PAGIENT_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:     "pagient",
		Version:  version.Version.String(),
		Usage:    "pagient server",
		Authors:  authors(cfg),
		Flags:    flags(cfg),
		Before:   before(cfg),
		Commands: command(cfg),
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

func authors(cfg *config.Config) []*cli.Author {
	return []*cli.Author{
		{
			Name:  "David Schneiderbauer",
			Email: "david.schneiderbauer@dschneiderbauer.me",
		},
	}
}

func flags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		// logging flags
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "set logging level",
			EnvVars:     []string{"TERRASTATE_LOG_LEVEL"},
			Destination: &cfg.Logs.Level,
		},
		&cli.BoolFlag{
			Name:        "log-colored",
			Value:       false,
			Usage:       "enable colored logging",
			EnvVars:     []string{"TERRASTATE_LOG_COLORED"},
			Destination: &cfg.Logs.Colored,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       false,
			Usage:       "enable pretty logging",
			EnvVars:     []string{"TERRASTATE_LOG_PRETTY"},
			Destination: &cfg.Logs.Pretty,
		},
	}
}

func before(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		switch strings.ToLower(cfg.Logs.Level) {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		if cfg.Logs.Pretty {
			log.Logger = log.Output(
				zerolog.ConsoleWriter{
					Out:     os.Stderr,
					NoColor: !cfg.Logs.Colored,
				},
			)
		}

		return nil
	}
}

func command(cfg *config.Config) []*cli.Command {
	return []*cli.Command{
		Server(cfg),
	}
}
