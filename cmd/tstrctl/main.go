package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tstrctl",
	Short: "tstrctl is the cli interface for tstr.",
	Long:  "tstrctl is the cli tool used to interacte and manage tstr.",
}

func init() {
	rootCmd.PersistentFlags().String("api-addr", "", "The address of the tstr gRPC API to dial.")
	viper.BindPFlag("api-addr", rootCmd.PersistentFlags().Lookup("api-addr"))

	rootCmd.PersistentFlags().Duration("timeout", 15*time.Second, "The amount of time to wait API requests.")
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
