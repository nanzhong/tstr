package main

import (
	"context"
	"fmt"

	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlTestGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a registered test",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withCtlControlClient(context.Background(), func(ctx context.Context, client controlv1.ControlServiceClient) error {
			fmt.Printf("Getting registered test %s...\n", args[0])

			ctx = metadata.AppendToOutgoingContext(ctx, auth.MDKeyNamespace, ctlTestNamespace)
			res, err := client.GetTest(ctx, &controlv1.GetTestRequest{Id: args[0]})
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	ctlTestCmd.AddCommand(ctlTestGetCmd)
}
