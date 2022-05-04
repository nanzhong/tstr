package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	var (
		file    string
		name    string
		labels  []string
		cron    string
		image   string
		command string
		args    []string
		env     []string
	)
	testRegisterCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new test.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			req := &control.RegisterTestRequest{
				RunConfig: &common.Test_RunConfig{},
			}

			if file != "" {
				testConfigBytes, err := os.ReadFile(file)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(testConfigBytes, &req); err != nil {
					return err
				}
			}

			if name != "" {
				req.Name = name
			}

			if len(labels) > 0 {
				for _, label := range labels {
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

			if cron != "" {
				req.CronSchedule = cron
			}

			if image != "" {
				req.RunConfig.ContainerImage = image
			}

			if command != "" {
				req.RunConfig.Command = command
			}

			if len(args) > 0 {
				req.RunConfig.Args = args
			}
			fmt.Println(args, req.RunConfig.Args)

			if len(env) > 0 {
				for _, ev := range env {
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

			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
			defer cancel()

			return withControlClient(ctx, func(ctx context.Context, client control.ControlServiceClient) error {
				res, err := client.RegisterTest(ctx, req)
				if err != nil {
					return err
				}

				fmt.Println(protojson.Format(res))
				return nil
			})
		},
	}

	testRegisterCmd.Flags().StringVar(&file, "file", "", "Path to a file containing the test definition.")
	testRegisterCmd.Flags().StringVar(&name, "name", "", "Name of the test.")
	testRegisterCmd.Flags().StringArrayVar(&labels, "labels", nil, "Labels for the test in key value pairs (e.g. key=value).")
	testRegisterCmd.Flags().StringVar(&cron, "cron", "", "Cron schedule the test should run at.")
	testRegisterCmd.Flags().StringVar(&image, "image", "", "Container image for the test.")
	testRegisterCmd.Flags().StringVar(&command, "command", "", "Command to run in the container.")
	testRegisterCmd.Flags().StringArrayVar(&args, "args", nil, "Arguments to pass to the command.")
	testRegisterCmd.Flags().StringArrayVar(&env, "env", nil, "Environment variable to execute the test with in key value pairs (e.g. key=value).")

	testCmd.AddCommand(testRegisterCmd)
}
