package main

import "github.com/spf13/cobra"

var ctlRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage test runs",
}

func init() {
	ctlTestCmd.PersistentFlags().StringVar(&ctlTestNamespace, "namespace", "", "The namespace to use.")
	ctlTestCmd.MarkPersistentFlagRequired("namespace")
	ctlCmd.AddCommand(ctlRunCmd)
}
