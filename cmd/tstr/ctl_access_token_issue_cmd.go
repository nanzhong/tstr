package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	ctlAccessTokenIssueName     string
	ctlAccessTokenValidDuration time.Duration
	ctlAccessTokenScopes        []string
	accessTokenIssueCmd         = &cobra.Command{
		Use:   "issue",
		Short: "Issue a new access token",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
			defer cancel()

			req := &adminv1.IssueAccessTokenRequest{
				Name: ctlAccessTokenIssueName,
			}

			for _, s := range ctlAccessTokenScopes {
				switch strings.ToLower(s) {
				case "admin":
					req.Scopes = append(req.Scopes, commonv1.AccessToken_ADMIN)
				case "control_r":
					req.Scopes = append(req.Scopes, commonv1.AccessToken_CONTROL_R)
				case "control_rw":
					req.Scopes = append(req.Scopes, commonv1.AccessToken_CONTROL_RW)
				case "runner":
					req.Scopes = append(req.Scopes, commonv1.AccessToken_RUNNER)
				case "data":
					req.Scopes = append(req.Scopes, commonv1.AccessToken_DATA)
				default:
					return fmt.Errorf("invalid access token scope %s", s)
				}
			}

			if ctlAccessTokenValidDuration != 0 {
				req.ValidDuration = durationpb.New(ctlAccessTokenValidDuration)
			}

			return withAdminClient(ctx, viper.GetString("ctl.grpc-addr"), !viper.GetBool("ctl.insecure"), viper.GetString("ctl.access-token"), func(ctx context.Context, client adminv1.AdminServiceClient) error {
				res, err := client.IssueAccessToken(ctx, req)
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
	accessTokenIssueCmd.Flags().StringVar(&ctlAccessTokenIssueName, "name", "", "The name for the access token.")
	accessTokenIssueCmd.MarkFlagRequired("name")
	accessTokenIssueCmd.Flags().DurationVar(&ctlAccessTokenValidDuration, "valid-duration", 7*24*time.Hour, "How long the token should be valid for (0 for non-expiring token).")
	accessTokenIssueCmd.Flags().StringArrayVar(&ctlAccessTokenScopes, "scopes", []string{"admin"}, "The scopes to attach to the access token.")

	accessTokenCmd.AddCommand(accessTokenIssueCmd)
}
