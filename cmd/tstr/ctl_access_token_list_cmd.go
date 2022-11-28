package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

type ctlAccessTokenListResult struct {
	res *adminv1.ListAccessTokensResponse
}

func (r *ctlAccessTokenListResult) RenderText(w io.Writer) error {
	columns := []table.Column{
		{Title: "ID", Width: 37},
		{Title: "Name", Width: 20},
		{Title: "Scopes", Width: 29},
		{Title: "NS Selectors", Width: 15},
		{Title: "Issued At (UTC)", Width: 20},
		{Title: "Expiry", Width: 10},
	}
	var rows []table.Row
	for _, t := range r.res.AccessTokens {
		var scopes []string
		for _, s := range t.Scopes {
			scopes = append(scopes, s.String())
		}

		row := table.Row{
			t.Id,
			t.Name,
			strings.Join(scopes, ", "),
			strings.Join(t.NamespaceSelectors, ", "),
			t.IssuedAt.AsTime().UTC().Format("2006-01-02 15:04:05"),
		}
		if t.ExpiresAt != nil {
			row = append(row, t.ExpiresAt.AsTime().Sub(t.IssuedAt.AsTime()).Truncate(time.Second).String())
		}
		rows = append(rows, row)
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Padding(0)
	s.Cell = s.Cell.Padding(0)
	s.Selected = s.Cell

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithStyles(s),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)

	_, err := fmt.Fprintln(w, fmt.Sprintf("Total # of access tokens: %d\n\n%s", len(r.res.AccessTokens), t.View()))
	return err
}

func (r *ctlAccessTokenListResult) RenderJSON(w io.Writer) error {
	_, err := fmt.Fprintln(w, protojson.Format(r.res))
	return err
}

var (
	ctlAccessTokenListIncludeExpired bool
	ctlAccessTokenListIncludeRevoked bool
	ctlAccessTokenListCmd            = &cobra.Command{
		Use:   "list",
		Short: "Retrieve information for all access tokens",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateOutputFormat(); err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
			defer cancel()

			req := &adminv1.ListAccessTokensRequest{
				IncludeExpired: ctlAccessTokenListIncludeRevoked,
				IncludeRevoked: ctlAccessTokenListIncludeRevoked,
			}

			return withAdminClient(ctx, viper.GetString("ctl.grpc-addr"), !viper.GetBool("ctl.insecure"), viper.GetString("ctl.access-token"), func(ctx context.Context, client adminv1.AdminServiceClient) error {
				res, err := client.ListAccessTokens(ctx, req)
				if err != nil {
					return err
				}

				return render(&ctlAccessTokenListResult{res: res}, cmd.OutOrStdout())
			})
		},
	}
)

func init() {
	ctlAccessTokenListCmd.Flags().BoolVar(&ctlAccessTokenListIncludeExpired, "include-expired", false, "Include expired tokens in result.")
	ctlAccessTokenListCmd.Flags().BoolVar(&ctlAccessTokenListIncludeRevoked, "include-revoked", false, "Include revoked tokens in result.")
	addOutputFlag(ctlAccessTokenListCmd)
	accessTokenCmd.AddCommand(ctlAccessTokenListCmd)
}
