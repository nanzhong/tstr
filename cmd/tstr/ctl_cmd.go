package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ctlCmd = &cobra.Command{
		Use:   "ctl",
		Short: "The cli interface for interacting with tstr",
	}
	ctlNamespace string
	ctlOutput    string
)

const (
	outputFormatText = "text"
	outputFormatJSON = "json"
)

func init() {
	ctlCmd.PersistentFlags().String("grpc-addr", "localhost:9000", "Address of the tstr gRPC API to dial.")
	viper.BindPFlag("ctl.grpc-addr", ctlCmd.PersistentFlags().Lookup("grpc-addr"))

	ctlCmd.PersistentFlags().Bool("insecure", false, "Insecure connection to api.")
	viper.BindPFlag("ctl.insecure", ctlCmd.PersistentFlags().Lookup("insecure"))

	ctlCmd.PersistentFlags().Duration("timeout", 15*time.Second, "Amount of time to wait API requests.")
	viper.BindPFlag("ctl.timeout", ctlCmd.PersistentFlags().Lookup("timeout"))

	ctlCmd.PersistentFlags().String("access-token", "", "Access token to use for authentication.")
	viper.BindPFlag("ctl.access-token", ctlCmd.PersistentFlags().Lookup("access-token"))

	rootCmd.AddCommand(ctlCmd)
}

func addOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ctlOutput, "output", "o", outputFormatText, "Output format [json, text]")
}

func validateOutputFormat() error {
	ctlOutput = strings.ToLower(ctlOutput)

	switch ctlOutput {
	case outputFormatText, outputFormatJSON:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", ctlOutput)
	}
}

func render(r outputRenderer, w io.Writer) error {
	switch ctlOutput {
	case outputFormatJSON:
		return r.RenderJSON(w)
	case outputFormatText:
		return r.RenderText(w)
	default:
		return fmt.Errorf("failed to  render: output format not specified")
	}
}

type outputRenderer interface {
	RenderText(io.Writer) error
	RenderJSON(io.Writer) error
}
