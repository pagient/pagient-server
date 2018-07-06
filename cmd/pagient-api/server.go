package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oklog/run"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/presenter/handler"
	"github.com/pagient/pagient-api/pkg/presenter/router"
	"github.com/pagient/pagient-api/pkg/presenter/websocket"
	"github.com/pagient/pagient-api/pkg/repository"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "server",
		Usage:  "start the integrated server",
		Before: serverBefore(cfg),
		Action: serverAction(cfg),
	}
}

func serverBefore(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		return nil
	}
}

func serverAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		// Initialize Repositories  (database access)
		clientRepo, err := repository.GetClientRepositoryInstance(cfg)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("client repository initialization failed")

			os.Exit(1)
		}
		pagerRepo, err := repository.GetPagerRepositoryInstance(cfg)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("pager repository initialization failed")

			os.Exit(1)
		}
		patientRepo, err := repository.GetPatientRepositoryInstance(cfg)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("patient repository initialization failed")

			os.Exit(1)
		}

		// Initialize Services (business logic)
		clientService := service.NewClientService(clientRepo)
		pagerService := service.NewPagerService(pagerRepo)
		patientService := service.NewPatientService(cfg, patientRepo, pagerRepo)

		// Initialize Websocket Hub and Handler (presenter layer)
		hub := websocket.NewHub()
		go hub.Run()

		clientHandler := handler.NewClientHandler(clientService)
		pagerHandler := handler.NewPagerHandler(pagerService)
		patientHandler := handler.NewPatientHandler(patientService, hub)
		websocketHandler := handler.NewWebsocketHandler(cfg, hub)

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

		if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				cfg.Server.Cert,
				cfg.Server.Key,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to load certificates")

				return err
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Address,
					Handler:      router.Load(cfg, clientHandler, pagerHandler, patientHandler, websocketHandler, clientService, patientService),
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
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Address,
				Handler:      router.Load(cfg, clientHandler, pagerHandler, patientHandler, websocketHandler, clientService, patientService),
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
					log.Info().
						Err(err).
						Msg("failed to stop http server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("http server stopped gracefully")
			})
		}

		return gr.Run()
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
