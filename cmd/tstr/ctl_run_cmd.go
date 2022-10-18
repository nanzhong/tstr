package main

import "github.com/spf13/cobra"

var ctlRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage test runs",
}

func init() {
	ctlRunCmd.PersistentFlags().StringVar(&ctlNamespace, "namespace", "", "The namespace to use.")
	ctlRunCmd.MarkPersistentFlagRequired("namespace")
	ctlCmd.AddCommand(ctlRunCmd)
}
