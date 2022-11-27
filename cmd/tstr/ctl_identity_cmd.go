package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	identityv1 "github.com/nanzhong/tstr/api/identity/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

type ctlIdentityResult struct {
	res *identityv1.IdentityResponse
}

func (r *ctlIdentityResult) RenderText(w io.Writer) error {
	var titleStyle = lipgloss.NewStyle().Bold(true).Width(20)

	fmt.Fprintf(w, "%s %s (%s)\n", titleStyle.Render("Access Token"), r.res.AccessToken.Name, r.res.AccessToken.Id)
	var scopes []string
	for _, s := range r.res.AccessToken.Scopes {
		scopes = append(scopes, s.String())
	}
	fmt.Fprintf(w, "%s %s\n", titleStyle.Render("Scopes"), strings.Join(scopes, ", "))
	fmt.Fprintf(w, "%s %s (accessible: %s)\n", titleStyle.Render("Namespace Selectors"), strings.Join(r.res.AccessToken.NamespaceSelectors, ", "), strings.Join(r.res.AccessibleNamespaces, ", "))
	fmt.Fprintf(w, "%s %s\n", titleStyle.Render("Issued At"), r.res.AccessToken.IssuedAt.AsTime().String())
	if !r.res.AccessToken.ExpiresAt.AsTime().IsZero() {
		fmt.Fprintf(w, "%s %s\n", titleStyle.Render("Expires At"), r.res.AccessToken.ExpiresAt.AsTime().String())
	}
	return nil
}

func (r *ctlIdentityResult) RenderJSON(w io.Writer) error {
	_, err := fmt.Fprintln(w, protojson.Format(r.res))
	return err
}

var ctlIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: "identity",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutputFormat(); err != nil {
			return err
		}

		return withCtlIdentityClient(cmd.Context(), func(ctx context.Context, isc identityv1.IdentityServiceClient) error {
			res, err := isc.Identity(ctx, &identityv1.IdentityRequest{})
			if err != nil {
				return err
			}

			return render(&ctlIdentityResult{res: res}, cmd.OutOrStdout())
		})
	},
}

func init() {
	addOutputFlag(ctlIdentityCmd)
	ctlCmd.AddCommand(ctlIdentityCmd)
}
