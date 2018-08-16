//go:generate fileb0x ../../b0x.yml

package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // import sqlite for database connection
	"github.com/oklog/run"
	"github.com/pagient/pagient-server/pkg/bridge"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/presenter/router"
	"github.com/pagient/pagient-server/pkg/presenter/websocket"
	"github.com/pagient/pagient-server/pkg/repository"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

// Server provides the sub-command to start the server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "start the integrated server",

		Before: func(c *cli.Context) error {
			return nil
		},

		Action: func(c *cli.Context) error {
			cfg, err := config.New()
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("config could not be loaded")

				os.Exit(1)
			}

			level, err := zerolog.ParseLevel(cfg.Log.Level)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("parse log level failed")
			}
			zerolog.SetGlobalLevel(level)

			logFile, err := os.OpenFile(path.Join(cfg.General.Root, "pagient.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("logfile could not be opened")

				os.Exit(1)
			}

			if cfg.Log.Pretty {
				log.Logger = log.Output(
					zerolog.ConsoleWriter{
						Out:     logFile,
						NoColor: !cfg.Log.Colored,
					},
				)
			} else {
				log.Logger = log.Output(logFile)
			}

			// Initialize Database Connection
			db, err := gorm.Open("sqlite3", cfg.Database.Path)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("establish database connection failed")

				os.Exit(1)
			}
			db.LogMode(zerolog.GlobalLevel() <= zerolog.DebugLevel)
			db.SetLogger(&log.Logger)
			defer db.Close()

			if err := repository.InitDatabase(db); err != nil {
				log.Fatal().
					Err(err).
					Msg("database initialization failed")

				os.Exit(1)
			}

			// Initialize Repositories  (database access)
			clientRepo, err := repository.GetClientRepositoryInstance(db)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("client repository initialization failed")

				os.Exit(1)
			}
			pagerRepo, err := repository.GetPagerRepositoryInstance(db)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("pager repository initialization failed")

				os.Exit(1)
			}
			patientRepo, err := repository.GetPatientRepositoryInstance(db)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("patient repository initialization failed")

				os.Exit(1)
			}
			tokenRepo, err := repository.GetTokenRepositoryInstance(db)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("token repository initialization failed")

				os.Exit(1)
			}
			userRepo, err := repository.GetUserRepositoryInstance(db)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("user repository initialization failed")

				os.Exit(1)
			}

			// Initialize Websocket Hub and Handler (presenter layer)
			hub := websocket.NewHub()

			// Initialize Services (business logic)
			clientService := service.NewClientService(clientRepo)
			pagerService := service.NewPagerService(pagerRepo)
			patientService := service.NewPatientService(cfg, patientRepo, pagerRepo, hub)
			tokenService := service.NewTokenService(cfg, tokenRepo)
			userService := service.NewUserService(cfg, userRepo)

			var gr run.Group

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)
					<-stop

					return nil
				}, func(err error) {
					close(stop)
				})
			}

			{
				stop := make(chan struct{}, 1)

				gr.Add(func() error {
					log.Info().
						Msg("starting websocket hub")

					hub.Run(stop)
					<-stop

					return nil
				}, func(reason error) {
					close(stop)
				})
			}

			{
				surgerySoftwareBridge := bridge.NewBridge(cfg, patientService, hub)
				stop := make(chan struct{}, 1)

				gr.Add(func() error {
					log.Info().
						Msg("starting surgery software bridge")

					return surgerySoftwareBridge.Run(stop)
				}, func(reason error) {
					close(stop)

					log.Info().
						AnErr("reason", reason).
						Msg("bridge stopped gracefully")
				})
			}

			if cfg.Server.Cert != "" && cfg.Server.Key != "" {
				cert, err := tls.LoadX509KeyPair(
					cfg.Server.Cert,
					cfg.Server.Key,
				)

				if err != nil {
					log.Fatal().
						Err(err).
						Msg("failed to load certificates")

					return err
				}

				{
					server := &http.Server{
						Addr:         cfg.Server.Address,
						Handler:      router.Load(cfg, clientService, pagerService, patientService, tokenService, userService, hub),
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 10 * time.Second,
						TLSConfig: &tls.Config{
							PreferServerCipherSuites: true,
							MinVersion:               tls.VersionTLS12,
							CurvePreferences:         curves(cfg),
							CipherSuites:             ciphers(cfg),
							Certificates:             []tls.Certificate{cert},
						},
					}

					gr.Add(func() error {
						log.Info().
							Str("addr", cfg.Server.Address).
							Msg("starting https server")

						return server.ListenAndServeTLS("", "")
					}, func(reason error) {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()

						if err := server.Shutdown(ctx); err != nil {
							log.Error().
								Err(err).
								Msg("failed to stop https server gracefully")

							return
						}

						log.Info().
							AnErr("reason", reason).
							Msg("https server stopped gracefully")
					})
				}

				return gr.Run()
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Address,
					Handler:      router.Load(cfg, clientService, pagerService, patientService, tokenService, userService, hub),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Address).
						Msg("starting http server")

					return server.ListenAndServe()
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Error().
							Err(err).
							Msg("failed to stop http server gracefully")

						return
					}

					log.Info().
						AnErr("reason", reason).
						Msg("http server stopped gracefully")
				})
			}

			return gr.Run()
		},
	}
}

func curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
