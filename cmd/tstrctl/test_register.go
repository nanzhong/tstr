package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new test.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		testConfigPath := args[0]
		testConfigBytes, err := os.ReadFile(testConfigPath)
		if err != nil {
			return err
		}

		var req *control.RegisterTestRequest
		if err := json.Unmarshal(testConfigBytes, &req); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()
		fmt.Println("dialing")
		return withControlClient(ctx, func(ctx context.Context, client control.ControlServiceClient) error {
			fmt.Println("starting")
			res, err := client.RegisterTest(ctx, req)
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(res)
		})
	},
}
