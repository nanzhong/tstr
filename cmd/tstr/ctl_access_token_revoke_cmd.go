package main

import (
	"context"
	"fmt"

	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlAccessTokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an access token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
		defer cancel()

		req := &adminv1.RevokeAccessTokenRequest{
			Id: args[0],
		}

		return withAdminClient(ctx, viper.GetString("ctl.grpc-addr"), !viper.GetBool("ctl.insecure"), viper.GetString("ctl.access-token"), func(ctx context.Context, client adminv1.AdminServiceClient) error {
			res, err := client.RevokeAccessToken(ctx, req)
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	accessTokenCmd.AddCommand(ctlAccessTokenRevokeCmd)
}
