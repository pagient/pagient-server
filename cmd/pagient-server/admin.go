package main

import (
	"os"

	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/database"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

// Admin provides the sub-command to perform administrative tasks
func Admin() *cli.Command {
	subcmdCreateUser := &cli.Command{
		Name: "create-user",
		Usage: "Create a new user in database",
		Action: runCreateUser,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Usage: "Username",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "User password",
			},
			&cli.UintFlag{
				Name: "client",
				Usage: "Client ID",
			},
		},
	}

	return &cli.Command{
		Name:  "admin",
		Usage: "perform admin specific tasks, e.g. create users and clients",
		Subcommands: []*cli.Command{
			subcmdCreateUser,
		},
	}
}

func basicSetup() service.Service {
	if err := config.Load(); err != nil {
		log.Fatal().
			Err(err).
			Msg("config could not be loaded")

		os.Exit(1)
	}

	// Setup Database Connection
	db, err := database.Open()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("database initialization failed")

		os.Exit(1)
	}

	// Setup Business Layer
	s := service.Init(db, nil)

	return s
}

func runCreateUser(c *cli.Context) error {
	s := basicSetup()

	user := &model.User{
		Username: c.String("username"),
		Password: c.String("password"),
		ClientID: c.Uint("client"),
	}

	user, err := s.CreateUser(user)
	if err != nil && service.IsModelValidationErr(err) {
		log.Info().
			Msgf("User is invalid: %s", err.Error())

		return nil
	}

	return errors.Wrap(err, "create user failed")
}
