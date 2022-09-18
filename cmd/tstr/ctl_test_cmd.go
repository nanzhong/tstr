package main

import "github.com/spf13/cobra"

var (
	ctlTestCmd = &cobra.Command{
		Use:   "test",
		Short: "test is the subcommand for interacting with test configuration",
	}
	ctlTestNamespace string
)

func init() {
	ctlTestCmd.PersistentFlags().StringVar(&ctlTestNamespace, "namespace", "", "The namespace to use.")
	ctlTestCmd.MarkFlagRequired("namespace")
	ctlCmd.AddCommand(ctlTestCmd)
}
