package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	ctlTestRegisterFile    string
	ctlTestRegisterID      string
	ctlTestRegisterName    string
	ctlTestRegisterLabels  []string
	ctlTestRegisterCron    string
	ctlTestRegisterImage   string
	ctlTestRegisterCommand string
	ctlTestRegisterArgs    []string
	ctlTestRegisterEnv     []string
	ctlTestRegisterTimeout time.Duration
	ctlTestRegisterCmd     = &cobra.Command{
		Use:   "register",
		Short: "Register a new test",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			test := &commonv1.Test{
				RunConfig: &commonv1.Test_RunConfig{},
			}

			if ctlTestRegisterFile != "" {
				testConfigBytes, err := os.ReadFile(ctlTestRegisterFile)
				if err != nil {
					return err
				}

				if err := protojson.Unmarshal(testConfigBytes, test); err != nil {
					return err
				}
			}

			var fieldMask fieldmaskpb.FieldMask

			if ctlTestRegisterID != "" {
				test.Id = ctlTestRegisterID
			}

			if ctlTestRegisterName != "" {
				fieldMask.Paths = append(fieldMask.Paths, "name")
				test.Name = ctlTestRegisterName
			}

			if len(ctlTestRegisterLabels) > 0 {
				fieldMask.Paths = append(fieldMask.Paths, "labels")
				for _, label := range ctlTestRegisterLabels {
					kv := strings.Split(label, "=")
					if len(kv) != 2 {
						continue
					}

					if test.Labels == nil {
						test.Labels = make(map[string]string)
					}
					test.Labels[kv[0]] = kv[1]
				}
			}

			if ctlTestRegisterCron != "" {
				fieldMask.Paths = append(fieldMask.Paths, "cron_schedule")
				test.CronSchedule = ctlTestRegisterCron
			}

			if ctlTestRegisterImage != "" {
				fieldMask.Paths = append(fieldMask.Paths, "run_config.container_image")
				test.RunConfig.ContainerImage = ctlTestRegisterImage
			}

			if ctlTestRegisterCommand != "" {
				fieldMask.Paths = append(fieldMask.Paths, "run_config.command")
				test.RunConfig.Command = ctlTestRegisterCommand
			}

			if len(ctlTestRegisterArgs) > 0 {
				fieldMask.Paths = append(fieldMask.Paths, "run_config.args")
				test.RunConfig.Args = ctlTestRegisterArgs
			}

			if len(ctlTestRegisterEnv) > 0 {
				fieldMask.Paths = append(fieldMask.Paths, "run_config.env")
				for _, ev := range ctlTestRegisterEnv {
					kv := strings.Split(ev, "=")
					if len(kv) != 2 {
						continue
					}

					if test.RunConfig.Env == nil {
						test.RunConfig.Env = make(map[string]string)
					}
					test.RunConfig.Env[kv[0]] = kv[1]
				}
			}

			if ctlTestRegisterTimeout > 0 {
				fieldMask.Paths = append(fieldMask.Paths, "run_config.timeout")
				test.RunConfig.Timeout = durationpb.New(ctlTestRegisterTimeout)
			}

			return withCtlControlClient(context.Background(), func(ctx context.Context, client controlv1.ControlServiceClient) error {
				if test.Id == "" {
					fmt.Printf("Registering new test: %s\n", test.Name)

					res, err := client.RegisterTest(ctx, &controlv1.RegisterTestRequest{
						Name:         test.Name,
						Labels:       test.Labels,
						RunConfig:    test.RunConfig,
						CronSchedule: test.CronSchedule,
					})
					if err != nil {
						return err
					}
					fmt.Println(protojson.Format(res))
				} else {
					fmt.Printf("Updating existing test: %s (%s)\n", test.Name, test.Id)

					res, err := client.UpdateTest(ctx, &controlv1.UpdateTestRequest{
						FieldMask:    &fieldMask,
						Id:           test.Id,
						Name:         test.Name,
						RunConfig:    test.RunConfig,
						Labels:       test.Labels,
						CronSchedule: test.CronSchedule,
					})
					if err != nil {
						return err
					}
					fmt.Println(protojson.Format(res))
				}
				return nil
			})
		},
	}
)

func init() {
	ctlTestRegisterCmd.Flags().StringVarP(&ctlTestRegisterFile, "file", "f", "", "Path to a file containing the test definition.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterID, "id", "", "ID of the test.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterName, "name", "", "Name of the test.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterLabels, "labels", nil, "Labels for the test in key value pairs (e.g. key=value).")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCron, "cron", "", "Cron schedule the test should run at.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterImage, "image", "", "Container image for the test.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCommand, "command", "", "Command to run in the container.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterArgs, "args", nil, "Arguments to pass to the command.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterEnv, "env", nil, "Environment variable to execute the test with in key value pairs (e.g. key=value).")
	ctlTestRegisterCmd.Flags().DurationVar(&ctlTestRegisterTimeout, "timeout", 5*time.Minute, "Timeout for the test execution.")

	ctlTestCmd.AddCommand(ctlTestRegisterCmd)
}
