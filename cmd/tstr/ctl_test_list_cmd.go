package main

import (
	"context"
	"fmt"

	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlTestListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered tests",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		return withCtlControlClient(context.Background(), func(ctx context.Context, client controlv1.ControlServiceClient) error {
			fmt.Println("Listing registered tests...")

			res, err := client.ListTests(ctx, &controlv1.ListTestsRequest{})
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	ctlTestCmd.AddCommand(ctlTestListCmd)
}
func init() {
	ctlCmd.AddCommand(ctlTestListCmd)
}
