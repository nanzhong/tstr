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
	accessTokenGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve information for an access token.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
			defer cancel()

			req := &admin.GetAccessTokenRequest{
				Id: args[0],
			}

			return withAdminClient(ctx, func(ctx context.Context, client admin.AdminServiceClient) error {
				res, err := client.GetAccessToken(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}

	accessTokenCmd.AddCommand(accessTokenGetCmd)
}
