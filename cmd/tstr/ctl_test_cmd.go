package main

import "github.com/spf13/cobra"

var ctlTestCmd = &cobra.Command{
	Use:   "test",
	Short: "test is the subcommand for interacting with test configuration",
}

func init() {
	ctlCmd.AddCommand(ctlTestCmd)
}
