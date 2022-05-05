package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
)

func init() {
	var (
		name          string
		validDuration time.Duration
		scopes        []string
	)

	accessTokenIssueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new access token.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
			defer cancel()

			req := &admin.IssueAccessTokenRequest{
				Name: name,
			}

			for _, s := range scopes {
				switch strings.ToLower(s) {
				case "admin":
					req.Scopes = append(req.Scopes, common.AccessToken_ADMIN)
				case "control_r":
					req.Scopes = append(req.Scopes, common.AccessToken_CONTROL_R)
				case "control_rw":
					req.Scopes = append(req.Scopes, common.AccessToken_CONTROL_RW)
				case "runner":
					req.Scopes = append(req.Scopes, common.AccessToken_RUNNER)
				default:
					return fmt.Errorf("invalid access token scope %s", s)
				}
			}

			if validDuration != 0 {
				req.ValidDuration = durationpb.New(validDuration)
			}

			return withAdminClient(ctx, func(ctx context.Context, client admin.AdminServiceClient) error {
				res, err := client.IssueAccessToken(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}

	accessTokenIssueCmd.Flags().StringVar(&name, "name", "", "The name for the access token.")
	accessTokenIssueCmd.MarkFlagRequired("name")
	accessTokenIssueCmd.Flags().DurationVar(&validDuration, "valid-duration", 7*24*time.Hour, "How long the token should be valid for (0 for non-expiring token).")
	accessTokenIssueCmd.Flags().StringArrayVar(&scopes, "scopes", []string{"admin"}, "The scopes to attach to the access token.")

	accessTokenCmd.AddCommand(accessTokenIssueCmd)
}
