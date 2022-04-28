package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tstrctl",
	Short: "tstrctl is the cli interface for tstr.",
	Long:  "tstrctl is the cli tool used to interacte and manage tstr.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
