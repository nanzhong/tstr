package main

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	grpczerolog "github.com/jwreagor/grpc-zerolog"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	datav1 "github.com/nanzhong/tstr/api/data/v1"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/server"
	"github.com/nanzhong/tstr/scheduler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "serve the gRPC and HTTP JSON API",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// TODO Use console writer for now for easy development/debugging, perhaps remove and rely on humanlog in the future.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		grpclog.SetLoggerV2(grpczerolog.New(log.Logger.With().Str("component", "client-grpc").Logger()))

		pgxPool, err := pgxpool.Connect(context.Background(), viper.GetString("api.pg-dsn"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to pg")
		}
		defer pgxPool.Close()

		if viper.GetString("api.bootstrap-token") != "" {
			tokenHashBytes := sha512.Sum512([]byte(viper.GetString("api.bootstrap-token")))
			tokenHash := hex.EncodeToString(tokenHashBytes[:])

			var textScopes []string
			for _, s := range []db.AccessTokenScope{db.AccessTokenScopeAdmin, db.AccessTokenScopeControlRw, db.AccessTokenScopeRunner, db.AccessTokenScopeData} {
				textScopes = append(textScopes, string(s))
			}
			_, err := db.New().IssueAccessToken(ctx, pgxPool, db.IssueAccessTokenParams{
				Name:      "bootstrap-token",
				TokenHash: tokenHash,
				Scopes:    textScopes,
				ExpiresAt: sql.NullTime{Valid: true, Time: time.Now().Add(24 * time.Hour)},
			})
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("failed to issue bootstrap token")
			}
		}

		l, err := net.Listen("tcp", viper.GetString("api.addr"))
		if err != nil {
			log.Fatal().
				Err(err).
				Str("addr", viper.GetString("api.addr")).
				Msg("failed to listen on api addr")
		}
		cm := cmux.New(l)
		grpcL := cm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		grpcgwL := cm.Match(cmux.HTTP1Fast())

		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				auth.UnaryServerInterceptor(pgxPool),
				grpc_validator.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				auth.StreamServerInterceptor(pgxPool),
				grpc_validator.StreamServerInterceptor(),
			),
		)

		scheduler := scheduler.New(pgxPool)

		controlServer := server.NewControlServer(pgxPool)
		controlv1.RegisterControlServiceServer(grpcServer, controlServer)

		adminServer := server.NewAdminServer(pgxPool)
		adminv1.RegisterAdminServiceServer(grpcServer, adminServer)

		runnerServer := server.NewRunnerServer(pgxPool)
		runnerv1.RegisterRunnerServiceServer(grpcServer, runnerServer)

		dataServer := server.NewDataServer(pgxPool)
		datav1.RegisterDataServiceServer(grpcServer, dataServer)

		grpcgwMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}))
		gwOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		datav1.RegisterDataServiceHandlerFromEndpoint(ctx, grpcgwMux, viper.GetString("api.addr"), gwOpts)
		grpcgwServer := http.Server{
			Handler: h2c.NewHandler(hlog.NewHandler(log.Logger)(grpcgwMux), &http2.Server{}),
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			<-ctx.Done()

			log.Info().Msg("shutting down")

			// Give one minute for running requests to complete
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			var eg errgroup.Group
			eg.Go(func() error {
				log.Info().Msg("attempting to shutdown scheduler")
				scheduler.Stop(shutdownCtx)
				return nil
			})
			eg.Go(func() error {
				log.Info().Msg("attempting to stop grpc server")
				grpcServer.GracefulStop()
				return nil
			})
			eg.Go(func() error {
				log.Info().Msg("attempting to shutdown grpc gw proxy")
				return grpcgwServer.Shutdown(shutdownCtx)
			})
			eg.Go(func() error {
				log.Info().Msg("attempting to close api mux")
				cm.Close()
				return nil
			})
			err := eg.Wait()
			if err != nil {
				log.Error().Err(err).Msg("failed to gracefully shutdown")
			}
		}()

		log.Info().Msg("tstr starting")
		var eg errgroup.Group
		eg.Go(func() error {
			log.Info().Msg("starting scheduler")
			return scheduler.Start()
		})
		eg.Go(func() error {
			log.Info().
				Msg("serving grpc server")
			return grpcServer.Serve(grpcL)
		})
		eg.Go(func() error {
			log.Info().
				Msg("serving grpc gw proxy")
			return grpcgwServer.Serve(grpcgwL)
		})
		eg.Go(func() error {
			log.Info().
				Str("addr", viper.GetString("api.addr")).
				Msg("serving api mux")
			return cm.Serve()
		})
		err = eg.Wait()
		log.Info().Msg("tstr shutdown")
	},
}

func init() {
	apiCmd.Flags().String("addr", "0.0.0.0:9000", "The address to serve the gRPC and HTTP JSON API on.")
	viper.BindPFlag("api.addr", apiCmd.Flags().Lookup("addr"))

	apiCmd.Flags().String("pg-dsn", "", "The PostgreSQL DSN to use.")
	viper.BindPFlag("api.pg-dsn", apiCmd.Flags().Lookup("pg-dsn"))

	apiCmd.Flags().String("bootstrap-token", "", "Bootstrap with provided access token (note that this token will have admin scope valid for 24h).")
	viper.BindPFlag("api.bootstrap-token", apiCmd.Flags().Lookup("bootstrap-token"))

	rootCmd.AddCommand(apiCmd)
}
