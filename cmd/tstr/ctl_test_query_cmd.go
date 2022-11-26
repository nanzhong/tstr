package main

import (
	"context"
	"fmt"

	datav1 "github.com/nanzhong/tstr/api/data/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlTestQueryCmd = &cobra.Command{
	Use:     "query",
	Short:   "query registered tests",
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"q", "list", "ls"},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return withCtlDataClient(context.Background(), func(ctx context.Context, client datav1.DataServiceClient) error {
			fmt.Println("Listing registered tests...")
			ctx = metadata.AppendToOutgoingContext(ctx, auth.MDKeyNamespace, ctlNamespace)
			res, err := client.QueryTests(ctx, &datav1.QueryTestsRequest{})
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	ctlTestCmd.AddCommand(ctlTestQueryCmd)
}
