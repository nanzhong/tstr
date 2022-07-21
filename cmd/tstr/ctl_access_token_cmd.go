package main

import "github.com/spf13/cobra"

var accessTokenCmd = &cobra.Command{
	Use:     "access-token",
	Short:   "Manage access tokens",
	Aliases: []string{"access", "access-tokens"},
}

func init() {
	ctlCmd.AddCommand(accessTokenCmd)
}
