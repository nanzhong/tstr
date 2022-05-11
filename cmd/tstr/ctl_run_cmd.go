package main

import "github.com/spf13/cobra"

var ctlRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage test runs.",
}

func init() {
	ctlCmd.AddCommand(ctlRunCmd)
}
