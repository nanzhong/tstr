package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tstr",
	Short: "tstr schedules, orchestrates, and reports on test runs.",
	Long:  "tstr is a tool that helps manage test workloads by taking care of test configuration, scheduling, orchestration, and reporting.",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
