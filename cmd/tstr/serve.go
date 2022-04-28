package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "sere the web UI",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Use console writer for now for easy development/debugging, perhaps remove and rely on humanlog in the future.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		l, err := net.Listen("tcp", viper.GetString("serve-api-addr"))
		if err != nil {
			log.Fatal().
				Err(err).
				Str("api-addre", viper.GetString("serve-api-addr")).
				Msg("failed to listen on api addr")
		}

		pool, err := pgxpool.Connect(context.Background(), viper.GetString("serve-pg-dsn"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to pg")
		}
		defer pool.Close()

		// TODO setup db.
		// TODO setup grpc server.
		// TODO setup http handler for web ui.
		mux := http.NewServeMux()
		httpServer := http.Server{
			Handler: mux,
		}

		// TODO setup startup/shutdown for grpc server, scheduler, etc.

		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			defer close(done)
			<-done

			log.Info().Msg("shutting down")
			{
				cancel()

				// Give one minute for running requests to complete
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
				defer cancel()

				var eg errgroup.Group
				eg.Go(func() error {
					log.Info().Msg("attempting to shutdown http server")
					return httpServer.Shutdown(shutdownCtx)
				})
				err := eg.Wait()
				if err != nil {
					log.Error().Err(err).Msg("failed to gracefully shutdown")
				}
			}
		}()

		var eg errgroup.Group
		eg.Go(func() error {
			log.Info().
				Str("web-addr", viper.GetString("serve-web-addr")).
				Msg("starting http server")
			return httpServer.Serve(l)
		})
		err = eg.Wait()
		log.Info().Msg("http server shutdown")
	},
}

func init() {
	serveCmd.Flags().String("api-addr", "0.0.0.0:9000", "The address to serve the gRPC API on.")
	viper.BindPFlag("serve-api-addr", serveCmd.Flags().Lookup("api-addr"))

	serveCmd.Flags().String("web-addr", "0.0.0.0:9090", "The address to serve the web UI on.")
	viper.BindPFlag("serve-web-addr", serveCmd.Flags().Lookup("web-addr"))

	serveCmd.Flags().String("pg-dsn", "", "The PostgreSQL DSN to use.")
	viper.BindPFlag("serve-pg-dsn", serveCmd.Flags().Lookup("pg-dsn"))
}
