package main

import (
	"context"
	"fmt"

	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ctlAccessTokenListIncludeExpired bool
	ctlAccessTokenListIncludeRevoked bool
	ctlAccessTokenListCmd            = &cobra.Command{
		Use:   "list",
		Short: "Retrieve information for all access tokens.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
			defer cancel()

			req := &admin.ListAccessTokensRequest{
				IncludeExpired: ctlAccessTokenListIncludeRevoked,
				IncludeRevoked: ctlAccessTokenListIncludeRevoked,
			}

			return withAdminClient(ctx, viper.GetString("ctl.api-addr"), viper.GetString("ctl.access-token"), func(ctx context.Context, client admin.AdminServiceClient) error {
				res, err := client.ListAccessTokens(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}
)

func init() {
	ctlAccessTokenListCmd.Flags().BoolVar(&ctlAccessTokenListIncludeExpired, "include-expired", false, "Include expired tokens in result.")
	ctlAccessTokenListCmd.Flags().BoolVar(&ctlAccessTokenListIncludeRevoked, "include-revoked", false, "Include revoked tokens in result.")

	accessTokenCmd.AddCommand(ctlAccessTokenListCmd)
}
