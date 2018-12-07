//go:generate fileb0x ../../b0x.yml

package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pagient/pagient-server/internal/bridge"
	bridgeDB "github.com/pagient/pagient-server/internal/bridge/database"
	"github.com/pagient/pagient-server/internal/caller"
	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/database"
	"github.com/pagient/pagient-server/internal/logger"
	"github.com/pagient/pagient-server/internal/presenter/router"
	"github.com/pagient/pagient-server/internal/presenter/websocket"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

// Web provides the sub-command to start the server.
func Web() *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "start the integrated webserver",

		Before: func(c *cli.Context) error {
			return nil
		},

		Action: func(c *cli.Context) error {
			if err := config.Load(); err != nil {
				log.Fatal().
					Err(err).
					Msg("config could not be loaded")

				os.Exit(1)
			}

			// Setup Logger
			if err := logger.Init(); err != nil {
				log.Fatal().
					Err(err).
					Msg("logger initialization failed")

				os.Exit(1)
			}
			defer logger.Close()

			// Setup Database Connection
			db, err := database.Open()
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("database initialization failed")

				os.Exit(1)
			}
			defer db.Close()

			// Initialize notifier (websocket hub)
			hub := websocket.NewHub()

			// Setup Business Layer
			s := service.Init(db, hub)

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
				// Setup Bridge Database Connection
				db, err := bridgeDB.Open()
				if err != nil {
					log.Fatal().
						Err(err).
						Msg("bridge database initialization failed")

					return err
				}

				// Setup Software Bridge
				softwareBridge := bridge.NewBridge(db)

				// Setup Caller
				pagerCaller := caller.NewCaller(s, softwareBridge)
				stop := make(chan struct{}, 1)

				gr.Add(func() error {
					log.Info().
						Msg("starting caller")

					return pagerCaller.Run(stop)
				}, func(reason error) {
					close(stop)

					log.Info().
						AnErr("reason", reason).
						Msg("caller stopped gracefully")
				})
			}

			if config.Server.Cert != "" && config.Server.Key != "" {
				cert, err := tls.LoadX509KeyPair(
					config.Server.Cert,
					config.Server.Key,
				)

				if err != nil {
					log.Fatal().
						Err(err).
						Msg("failed to load certificates")

					return err
				}

				{
					server := &http.Server{
						Addr:         config.Server.Address,
						Handler:      router.Load(s, hub),
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 10 * time.Second,
						TLSConfig: &tls.Config{
							PreferServerCipherSuites: true,
							MinVersion:               tls.VersionTLS12,
							CurvePreferences:         curves(),
							CipherSuites:             ciphers(),
							Certificates:             []tls.Certificate{cert},
						},
					}

					gr.Add(func() error {
						log.Info().
							Str("addr", config.Server.Address).
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
					Addr:         config.Server.Address,
					Handler:      router.Load(s, hub),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", config.Server.Address).
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

func curves() []tls.CurveID {
	if config.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers() []uint16 {
	if config.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
