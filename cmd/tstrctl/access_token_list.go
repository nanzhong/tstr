package main

import (
	"context"
	"fmt"

	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	var (
		includeExpired bool
		includeRevoked bool
	)

	accessTokenListCmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieve information for all access tokens.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
			defer cancel()

			req := &admin.ListAccessTokensRequest{
				IncludeExpired: includeExpired,
				IncludeRevoked: includeRevoked,
			}

			return withAdminClient(ctx, func(ctx context.Context, client admin.AdminServiceClient) error {
				res, err := client.ListAccessTokens(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}

	accessTokenListCmd.Flags().BoolVar(&includeExpired, "include-expired", false, "Include expired tokens in result.")
	accessTokenListCmd.Flags().BoolVar(&includeRevoked, "include-revoked", false, "Include revoked tokens in result.")

	accessTokenCmd.AddCommand(accessTokenListCmd)
}
