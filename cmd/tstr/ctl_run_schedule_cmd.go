package main

import (
	"context"
	"fmt"

	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var ctlRunScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a test run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
		defer cancel()

		return withControlClient(ctx, viper.GetString("ctl.api-addr"), viper.GetBool("ctl.secure"), viper.GetString("ctl.access-token"), func(ctx context.Context, client controlv1.ControlServiceClient) error {
			res, err := client.ScheduleRun(ctx, &controlv1.ScheduleRunRequest{
				TestId: args[0],
			})
			if err != nil {
				return err
			}

			fmt.Println(protojson.Format(res))
			return nil
		})
	},
}

func init() {
	ctlRunCmd.AddCommand(ctlRunScheduleCmd)
}
