package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/runner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a tstr runner for executing test workloads",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		var (
			acceptLabelSelectors = make(map[string]string)
			rejectLabelSelectors = make(map[string]string)
		)
		for _, l := range viper.GetStringSlice("run.accept-label-selectors") {
			parts := strings.Split(l, "=")
			if len(parts) != 2 {
				log.Fatal().Msg("invalid accept label selectors")
			}
			acceptLabelSelectors[parts[0]] = parts[1]
		}
		for _, l := range viper.GetStringSlice("run.reject-label-selectors") {
			parts := strings.Split(l, "=")
			if len(parts) != 2 {
				log.Fatal().Msg("invalid reject label selectors")
			}
			rejectLabelSelectors[parts[0]] = parts[1]
		}

		withRunnerClient(context.Background(), viper.GetString("run.api-addr"), viper.GetString("run.access-token"), func(ctx context.Context, client runnerv1.RunnerServiceClient) error {
			runner, err := runner.New(
				client,
				viper.GetString("run.name"),
				acceptLabelSelectors,
				rejectLabelSelectors,
			)
			if err != nil {
				return fmt.Errorf("failed to initialize runner: %w", err)
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			go func() {
				<-ctx.Done()
				ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("run.graceful-shutdown"))
				defer cancel()
				log.Info().Msg("stopping runner")
				runner.Stop(ctx)
			}()

			log.Info().
				Str("name", viper.GetString("run.name")).
				Msg("starting runner")
			if err := runner.Run(); err != nil {
				log.Error().Err(err).Msg("runner shutdown with error")
				return err
			}

			log.Info().Msg("runner shutdown gracefully")
			return nil
		})
	},
}

func init() {
	runCmd.Flags().String("api-addr", "0.0.0.0:9000", "Address of the tstr api to dial.")
	viper.BindPFlag("run.api-addr", runCmd.Flags().Lookup("api-addr"))

	runCmd.Flags().String("access-token", "", "Runner access token to use.")
	viper.BindPFlag("run.access-token", runCmd.Flags().Lookup("access-token"))

	hostname, _ := os.Hostname()
	runCmd.Flags().String("name", hostname, "Name of the runner.")
	viper.BindPFlag("run.name", runCmd.Flags().Lookup("name"))

	runCmd.Flags().StringArray("accept-label-selectors", nil, "Label selectors for test to accept.")
	viper.BindPFlag("run.accept-label-selectors", runCmd.Flags().Lookup("accept-label-selectors"))

	runCmd.Flags().StringArray("reject-label-selectors", nil, "Label selectors for test to reject.")
	viper.BindPFlag("run.reject-label-selectors", runCmd.Flags().Lookup("reject-label-selectors"))

	runCmd.Flags().Duration("graceful-shutdown", 5*time.Minute, "Amount of time to allow for graceful shutdown.")
	viper.BindPFlag("run.graceful-shutdown", runCmd.Flags().Lookup("graceful-shutdown"))

	rootCmd.AddCommand(runCmd)
}
