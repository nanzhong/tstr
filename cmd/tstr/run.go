package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	runnerapi "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/runner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Start a tstr runner for executing test workloads.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			conn, err := grpc.Dial(
				viper.GetString("run-api-addr"),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithChainUnaryInterceptor(
					grpc_validator.UnaryClientInterceptor(),
					auth.UnaryClientInterceptor(viper.GetString("run-access-token")),
				),
				grpc.WithChainStreamInterceptor(
					auth.StreamClientInterceptor(viper.GetString("run-access-token")),
				),
			)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to dial tstr api")
			}
			defer conn.Close()

			client := runnerapi.NewRunnerServiceClient(conn)
			runner := runner.New(
				client,
				viper.GetString("run-name"),
				nil,
				nil,
			)

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			go func() {
				<-ctx.Done()
				ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("run-graceful-shutdown"))
				defer cancel()
				log.Info().Msg("stopping runner")
				runner.Stop(ctx)
			}()

			log.Info().
				Str("name", viper.GetString("run-name")).
				Msg("starting runner")
			if err := runner.Run(); err != nil {
				log.Error().Err(err).Msg("runner shutdown with error")
				return
			}
			log.Info().Msg("runner shutdown gracefully")
		},
	}

	runCmd.Flags().String("api-addr", "", "Address of the tstr api to dial.")
	viper.BindPFlag("run-api-addr", runCmd.Flags().Lookup("api-addr"))

	runCmd.Flags().String("access-token", "", "Runner access token to use.")
	viper.BindPFlag("run-access-token", runCmd.Flags().Lookup("access-token"))

	hostname, _ := os.Hostname()
	runCmd.Flags().String("name", hostname, "Name of the runner.")
	viper.BindPFlag("run-name", runCmd.Flags().Lookup("name"))

	runCmd.Flags().StringArray("accept-label-selectors", nil, "Label selectors for test to accept.")
	viper.BindPFlag("run-accept-label-selectors", runCmd.Flags().Lookup("accept-label-selectors"))

	runCmd.Flags().StringArray("reject-label-selectors", nil, "Label selectors for test to reject.")
	viper.BindPFlag("run-reject-label-selectors", runCmd.Flags().Lookup("reject-label-selectors"))

	runCmd.Flags().Duration("graceful-shutdown", 5*time.Minute, "Amount of time to allow for graceful shutdown.")
	viper.BindPFlag("run-graceful-shutdown", runCmd.Flags().Lookup("graceful-shutdown"))

	rootCmd.AddCommand(runCmd)
}
