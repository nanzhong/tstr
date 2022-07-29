package main

import (
	"context"
	"crypto/x509"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func clientDialOpts(secure bool, accessToken string) []grpc.DialOption {
	var opts []grpc.DialOption
	if secure {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(
		opts,
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			grpc_validator.UnaryClientInterceptor(),
			auth.UnaryClientInterceptor(accessToken),
		),
		grpc.WithChainStreamInterceptor(
			auth.StreamClientInterceptor(accessToken),
		),
	)
	return opts
}

func withControlClient(ctx context.Context, apiAddr string, secure bool, accessToken string, fn func(context.Context, controlv1.ControlServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(secure, accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := controlv1.NewControlServiceClient(conn)
	return fn(ctx, client)
}

func withAdminClient(ctx context.Context, apiAddr string, secure bool, accessToken string, fn func(context.Context, adminv1.AdminServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(secure, accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := adminv1.NewAdminServiceClient(conn)
	return fn(ctx, client)
}

func withRunnerClient(ctx context.Context, apiAddr string, secure bool, accessToken string, fn func(context.Context, runnerv1.RunnerServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(secure, accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := runnerv1.NewRunnerServiceClient(conn)
	return fn(ctx, client)
}
