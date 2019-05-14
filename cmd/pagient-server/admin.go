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
		Name:   "create-user",
		Usage:  "Create a new user in database",
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
				Name:  "client",
				Usage: "Client ID",
			},
		},
	}

	subcmdChangePassword := &cli.Command{
		Name:   "change-password",
		Usage:  "Change a user's password",
		Action: runChangePassword,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Usage: "The user to change password for",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "New password to set for user",
			},
		},
	}

	subcmdCreateClient := &cli.Command{
		Name:   "create-client",
		Usage:  "Create a new client in database",
		Action: runCreateClient,
		Flags:  []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name",
			},
		},
	}

	subcmdCreatePager := &cli.Command{
		Name:   "create-pager",
		Usage:  "Create a new pager in database",
		Action: runCreatePager,
		Flags:  []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name",
			},
			&cli.UintFlag{
				Name:  "id",
				Usage: "EasyCall ID",
			},
		},
	}

	return &cli.Command{
		Name:  "admin",
		Usage: "perform admin specific tasks, e.g. create users and clients",
		Subcommands: []*cli.Command{
			subcmdCreateUser,
			subcmdChangePassword,
			subcmdCreateClient,
			subcmdCreatePager,
		},
	}
}

type db interface {
	Close() error
}

func basicSetup() (service.Service, db) {
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
	s := service.NewService(db, nil)

	return s, db
}

func runCreateUser(c *cli.Context) error {
	s, db := basicSetup()
	defer db.Close()

	user := &model.User{
		Username: c.String("username"),
		Password: c.String("password"),
		ClientID: c.Uint("client"),
	}

	err := s.CreateUser(user)
	if err != nil && service.IsModelValidationErr(err) {
		log.Info().
			Msgf("User is invalid: %s", err.Error())

		return nil
	}

	return errors.Wrap(err, "create user failed")
}

func runChangePassword(c *cli.Context) error {
	s, db := basicSetup()
	defer db.Close()

	user := &model.User{
		Username: c.String("username"),
		Password: c.String("password"),
	}

	err := s.ChangeUserPassword(user)
	if err != nil && service.IsModelValidationErr(err) {
		log.Info().
			Msgf("User is invalid: %s", err.Error())

		return nil
	}

	return errors.Wrap(err, "change user password failed")
}

func runCreateClient(c *cli.Context) error {
	s, db := basicSetup()
	defer db.Close()

	client := &model.Client{
		Name: c.String("name"),
	}

	err := s.CreateClient(client)
	if err != nil && service.IsModelValidationErr(err) {
		log.Info().
			Msgf("Client is invalid: %s", err.Error())

		return nil
	}

	return errors.Wrap(err, "create client failed")
}

func runCreatePager(c *cli.Context) error {
	s, db := basicSetup()
	defer db.Close()

	pager := &model.Pager{
		Name: c.String("name"),
		EasyCallID: c.Uint("id"),
	}

	err := s.CreatePager(pager)
	if err != nil && service.IsModelValidationErr(err) {
		log.Info().
			Msgf("Pager is invalid: %s", err.Error())

		return nil
	}

	return errors.Wrap(err, "create pager failed")
}
