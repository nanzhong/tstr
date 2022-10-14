package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	ctlTestRegisterFile         string
	ctlTestRegisterID           string
	ctlTestRegisterName         string
	ctlTestRegisterLabels       []string
	ctlTestRegisterMatrixLabels []string
	ctlTestRegisterCron         string
	ctlTestRegisterImage        string
	ctlTestRegisterCommand      string
	ctlTestRegisterArgs         []string
	ctlTestRegisterEnv          []string
	ctlTestRegisterTimeout      time.Duration
	ctlTestRegisterCmd          = &cobra.Command{
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

			if ctlTestRegisterID != "" {
				test.Id = ctlTestRegisterID
			}

			if ctlTestRegisterName != "" {
				test.Name = ctlTestRegisterName
			}

			if len(ctlTestRegisterLabels) > 0 {
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

			if len(ctlTestRegisterMatrixLabels) > 0 {
				for _, label := range ctlTestRegisterMatrixLabels {
					kv := strings.Split(label, "=")
					if len(kv) != 2 {
						continue
					}

					if test.Matrix == nil {
						test.Matrix = &commonv1.Test_Matrix{}
					}
					if test.Matrix.Labels == nil {
						test.Matrix.Labels = make(map[string]*commonv1.Test_Matrix_LabelValues)
					}
					if test.Matrix.Labels[kv[0]] == nil {
						test.Matrix.Labels[kv[0]] = &commonv1.Test_Matrix_LabelValues{}
					}
					test.Matrix.Labels[kv[0]].Values = append(test.Matrix.Labels[kv[0]].Values, kv[1])
				}
			}

			if ctlTestRegisterCron != "" {
				test.CronSchedule = ctlTestRegisterCron
			}

			if ctlTestRegisterImage != "" {
				test.RunConfig.ContainerImage = ctlTestRegisterImage
			}

			if ctlTestRegisterCommand != "" {
				test.RunConfig.Command = ctlTestRegisterCommand
			}

			if len(ctlTestRegisterArgs) > 0 {
				test.RunConfig.Args = ctlTestRegisterArgs
			}

			if len(ctlTestRegisterEnv) > 0 {
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
				test.RunConfig.Timeout = durationpb.New(ctlTestRegisterTimeout)
			}

			return withCtlControlClient(context.Background(), func(ctx context.Context, client controlv1.ControlServiceClient) error {
				ctx = metadata.AppendToOutgoingContext(ctx, auth.MDKeyNamespace, ctlTestNamespace)
				if test.Id == "" {
					fmt.Printf("Registering new test: %s\n", test.Name)

					ctx = metadata.AppendToOutgoingContext(ctx, auth.MDKeyNamespace, ctlTestNamespace)
					res, err := client.RegisterTest(ctx, &controlv1.RegisterTestRequest{
						Name:         test.Name,
						Labels:       test.Labels,
						Matrix:       test.Matrix,
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
						FieldMask: &fieldmaskpb.FieldMask{
							Paths: []string{
								"name",
								"labels",
								"matrix",
								"run_config.container_image",
								"run_config.command",
								"run_config.args",
								"run_config.env",
								"run_config.timeout",
								"cron_schedule",
							},
						},
						Id:           test.Id,
						Name:         test.Name,
						Labels:       test.Labels,
						Matrix:       test.Matrix,
						RunConfig:    test.RunConfig,
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
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterMatrixLabels, "matrix-labels", nil, "Test matrix label set (e.g. key=value).")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCron, "cron", "", "Cron schedule the test should run at.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterImage, "image", "", "Container image for the test.")
	ctlTestRegisterCmd.Flags().StringVar(&ctlTestRegisterCommand, "command", "", "Command to run in the container.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterArgs, "args", nil, "Arguments to pass to the command.")
	ctlTestRegisterCmd.Flags().StringArrayVar(&ctlTestRegisterEnv, "env", nil, "Environment variable to execute the test with in key value pairs (e.g. key=value).")
	ctlTestRegisterCmd.Flags().DurationVar(&ctlTestRegisterTimeout, "timeout", 5*time.Minute, "Timeout for the test execution.")

	ctlTestCmd.AddCommand(ctlTestRegisterCmd)
}
