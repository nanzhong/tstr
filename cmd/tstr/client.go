package main

import (
	"context"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientDialOpts(accessToken string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			grpc_validator.UnaryClientInterceptor(),
			auth.UnaryClientInterceptor(accessToken),
		),
		grpc.WithChainStreamInterceptor(
			auth.StreamClientInterceptor(accessToken),
		),
	}
}

func withControlClient(ctx context.Context, apiAddr string, accessToken string, fn func(context.Context, control.ControlServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := control.NewControlServiceClient(conn)
	return fn(ctx, client)
}

func withAdminClient(ctx context.Context, apiAddr string, accessToken string, fn func(context.Context, admin.AdminServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewAdminServiceClient(conn)
	return fn(ctx, client)
}

func withRunnerClient(ctx context.Context, apiAddr string, accessToken string, fn func(context.Context, runner.RunnerServiceClient) error) error {
	conn, err := grpc.Dial(
		apiAddr,
		clientDialOpts(accessToken)...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := runner.NewRunnerServiceClient(conn)
	return fn(ctx, client)
}
