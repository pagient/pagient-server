package main

import (
	"fmt"
	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/database"
	"github.com/pagient/pagient-server/internal/logger"
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
		Action: cliEnvSetup(runCreateUser),
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
		Action: cliEnvSetup(runChangePassword),
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
		Action: cliEnvSetup(runCreateClient),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name",
			},
		},
	}

	subcmdCreatePager := &cli.Command{
		Name:   "create-pager",
		Usage:  "Create a new pager in database",
		Action: cliEnvSetup(runCreatePager),
		Flags: []cli.Flag{
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

type commandFunc func(*cli.Context, service.Service, database.DB) error

func cliEnvSetup(cmdFunc commandFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := config.Load(); err != nil {
			log.Fatal().
				Err(err).
				Msg("config could not be loaded")
		}

		config.Log.Pretty = true

		// Setup Logger
		if err := logger.Init(); err != nil {
			log.Fatal().
				Err(err).
				Msg("logger initialization failed")
		}
		defer logger.Close()

		// Setup Database Connection
		db, err := database.Open()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("database initialization failed")
		}
		defer db.Close()

		// Setup Business Layer
		s := service.NewService(db, nil)

		return cmdFunc(c, s, db)
	}
}

func runCreateUser(c *cli.Context, s service.Service, db database.DB) error {
	user := &model.User{
		Username: c.String("username"),
		Password: c.String("password"),
		ClientID: c.Uint("client"),
	}

	err := s.CreateUser(user)
	if err != nil && service.IsModelValidationErr(err) {
		fmt.Printf("User is invalid: %s\n", err.Error())
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "create user failed")
	}

	fmt.Printf("User - ID %d - successfully created!\n", user.ID)
	return nil
}

func runChangePassword(c *cli.Context, s service.Service, db database.DB) error {
	user := &model.User{
		Username: c.String("username"),
		Password: c.String("password"),
	}

	err := s.ChangeUserPassword(user)
	if err != nil && service.IsModelValidationErr(err) {
		fmt.Printf("User is invalid: %s\n", err.Error())
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "change user password failed")
	}

	fmt.Printf("Password of User %s successfully changed!\n", user.Username)
	return nil
}

func runCreateClient(c *cli.Context, s service.Service, db database.DB) error {
	client := &model.Client{
		Name: c.String("name"),
	}

	err := s.CreateClient(client)
	if err != nil && service.IsModelValidationErr(err) {
		fmt.Printf("Client is invalid: %s\n", err.Error())
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "create client failed")
	}

	fmt.Printf("Client - ID %d - successfully created!\n", client.ID)
	return nil
}

func runCreatePager(c *cli.Context, s service.Service, db database.DB) error {
	pager := &model.Pager{
		Name:       c.String("name"),
		EasyCallID: c.Uint("id"),
	}

	err := s.CreatePager(pager)
	if err != nil && service.IsModelValidationErr(err) {
		fmt.Printf("Pager is invalid: %s\n", err.Error())
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "create pager failed")
	}

	fmt.Printf("Pager - ID %d - successfully created!\n", pager.ID)
	return nil
}
