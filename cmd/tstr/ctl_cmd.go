package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ctlCmd = &cobra.Command{
	Use:   "ctl",
	Short: "The cli interface for interacting with tstr",
}

func init() {
	ctlCmd.PersistentFlags().String("api-addr", "", "Address of the tstr gRPC API to dial.")
	viper.BindPFlag("ctl.api-addr", ctlCmd.PersistentFlags().Lookup("api-addr"))

	ctlCmd.PersistentFlags().Duration("timeout", 15*time.Second, "Amount of time to wait API requests.")
	viper.BindPFlag("ctl.timeout", ctlCmd.PersistentFlags().Lookup("timeout"))

	ctlCmd.PersistentFlags().String("access-token", "", "Access token to use for authentication.")
	viper.BindPFlag("ctl.access-token", ctlCmd.PersistentFlags().Lookup("access-token"))

	rootCmd.AddCommand(ctlCmd)
}
