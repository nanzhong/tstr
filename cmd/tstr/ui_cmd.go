package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nanzhong/tstr/ui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Serve the tstr web ui",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// TODO Use console writer for now for easy development/debugging, perhaps remove and rely on humanlog in the future.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		apiURL, err := url.Parse(viper.GetString("ui.api-base-url"))
		if err != nil {
			log.Fatal().Err(err).Str("api_base_url", viper.GetString("api_base_url")).Msg("invalid api base url")
		}

		uiProxy := ui.NewAPIProxy(apiURL, viper.GetString("ui.access-token"))
		uiServer := ui.NewUIServer()

		mux := http.NewServeMux()
		mux.Handle("/api/", uiProxy)
		mux.Handle("/", uiServer)

		server := http.Server{
			Addr:    viper.GetString("ui.addr"),
			Handler: mux,
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			<-ctx.Done()

			log.Info().Msg("attempting graceful shutdown")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("error shutting down http server")
			}
		}()

		log.Info().Msg("starting")

		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("stopped listening")
		}

		log.Info().Msg("shutdown")
	},
}

func init() {
	uiCmd.Flags().String("addr", "0.0.0.0:8000", "The address to serve the web ui on")
	viper.BindPFlag("ui.addr", uiCmd.Flags().Lookup("addr"))

	uiCmd.Flags().String("api-base-url", "http://0.0.0.0:9000", "The base URL of the tstr api to connect to")
	viper.BindPFlag("ui.api-base-url", uiCmd.Flags().Lookup("api-base-url"))

	uiCmd.Flags().String("access-token", "", "The access token to use when communicating with the tstr api (it must have the data scope)")
	viper.BindPFlag("ui.access-token", uiCmd.Flags().Lookup("access-token"))

	rootCmd.AddCommand(uiCmd)
}
