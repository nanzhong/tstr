package main

import (
	"context"
	"fmt"

	identityv1 "github.com/nanzhong/tstr/api/identity/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: "identity",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withCtlIdentityClient(cmd.Context(), func(ctx context.Context, isc identityv1.IdentityServiceClient) error {
			res, err := isc.Identity(ctx, &identityv1.IdentityRequest{})
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	ctlCmd.AddCommand(ctlIdentityCmd)
}
