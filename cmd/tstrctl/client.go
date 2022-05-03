package main

import (
	"context"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithBlock(),
	grpc.WithChainUnaryInterceptor(
		grpc_validator.UnaryClientInterceptor(),
	),
}

func withControlClient(ctx context.Context, fn func(context.Context, control.ControlServiceClient) error) error {
	conn, err := grpc.Dial(
		viper.GetString("api-addr"),
		clientDialOpts...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := control.NewControlServiceClient(conn)
	return fn(ctx, client)
}

func withAdminClient(ctx context.Context, fn func(context.Context, admin.AdminServiceClient) error) error {
	conn, err := grpc.Dial(
		viper.GetString("api-addr"),
		clientDialOpts...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewAdminServiceClient(conn)
	return fn(ctx, client)
}
