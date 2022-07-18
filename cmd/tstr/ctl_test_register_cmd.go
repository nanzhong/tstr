package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ctlTestRegisterFile    string
	ctlTestRegisterName    string
	ctlTestRegisterLabels  []string
	ctlTestRegisterCron    string
	ctlTestRegisterImage   string
	ctlTestRegisterCommand string
	ctlTestRegisterArgs    []string
	ctlTestRegisterEnv     []string
	ctlTestRegisterCmd     = &cobra.Command{
		Use:   "register",
		Short: "Register a new test.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			req := &controlv1.RegisterTestRequest{
				RunConfig: &commonv1.Test_RunConfig{},
			}

			if ctlTestRegisterFile != "" {
				testConfigBytes, err := os.ReadFile(ctlTestRegisterFile)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(testConfigBytes, &req); err != nil {
					return err
				}
			}

			if ctlTestRegisterName != "" {
				req.Name = ctlTestRegisterName
			}

			if len(ctlTestRegisterLabels) > 0 {
				for _, label := range ctlTestRegisterLabels {
					kv := strings.Split(label, "=")
					if len(kv) != 2 {
						continue
					}

					if req.Labels == nil {
						req.Labels = make(map[string]string)
					}
					req.Labels[kv[0]] = kv[1]
				}
			}

			if ctlTestRegisterCron != "" {
				req.CronSchedule = ctlTestRegisterCron
			}

			if ctlTestRegisterImage != "" {
				req.RunConfig.ContainerImage = ctlTestRegisterImage
			}

			if ctlTestRegisterCommand != "" {
				req.RunConfig.Command = ctlTestRegisterCommand
			}

			if len(ctlTestRegisterArgs) > 0 {
				req.RunConfig.Args = ctlTestRegisterArgs
			}

			if len(ctlTestRegisterEnv) > 0 {
				for _, ev := range ctlTestRegisterEnv {
					kv := strings.Split(ev, "=")
					if len(kv) != 2 {
						continue
					}

					if req.RunConfig.Env == nil {
						req.RunConfig.Env = make(map[string]string)
					}
					req.RunConfig.Env[kv[0]] = kv[1]
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
			defer cancel()

			return withControlClient(ctx, viper.GetString("ctl.api-addr"), viper.GetString("ctl.access-token"), func(ctx context.Context, client controlv1.ControlServiceClient) error {
				res, err := client.RegisterTest(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}
)

func init() {
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterFile, "file", "", "Path to a file containing the test definition.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterName, "name", "", "Name of the test.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterLabels, "labels", nil, "Labels for the test in key value pairs (e.g. key=value).")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCron, "cron", "", "Cron schedule the test should run at.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterImage, "image", "", "Container image for the test.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCommand, "command", "", "Command to run in the container.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterArgs, "args", nil, "Arguments to pass to the command.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterEnv, "env", nil, "Environment variable to execute the test with in key value pairs (e.g. key=value).")

	ctlTestCmd.AddCommand(ctlTestRegisterCmd)
}
