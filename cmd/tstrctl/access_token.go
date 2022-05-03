package main

import "github.com/spf13/cobra"

var accessTokenCmd = &cobra.Command{
	Use:     "access-token",
	Short:   "access-token is the subcommand for managing access tokens.",
	Aliases: []string{"access", "access-tokens"},
}

func init() {
	rootCmd.AddCommand(accessTokenCmd)
}
