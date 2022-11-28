package main

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type ctlRunScheduleResult struct {
	testID string
	res    *controlv1.ScheduleRunResponse
}

func (r *ctlRunScheduleResult) RenderText(w io.Writer) error {
	bold := lipgloss.NewStyle().Bold(true)

	if len(r.res.Runs) == 0 {
		_, err := fmt.Fprintf(w, "No runs scheduled for test %s", bold.Render(r.testID))
		return err
	}

	var output []string
	output = append(output, fmt.Sprintf("Scheduled %s runs for test %s", bold.Render(strconv.Itoa(len(r.res.Runs))), bold.Render(r.testID)))
	if r.res.Runs[0].TestMatrixId != "" {
		output = append(output, fmt.Sprintf("Test matrix ID for runs: %s\n", bold.Render(r.res.Runs[0].TestMatrixId)))
	}

	cols := []table.Column{
		{Title: "Run ID", Width: 37},
		{Title: "Labels", Width: 55},
	}
	var rows []table.Row
	for _, r := range r.res.Runs {
		var labels []string
		for k, v := range r.Labels {
			labels = append(labels, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(labels)
		rows = append(rows, table.Row{r.Id, strings.Join(labels, ", ")})
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Padding(0)
	s.Cell = s.Cell.Padding(0)
	s.Selected = s.Cell

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithStyles(s),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)
	output = append(output, t.View())

	_, err := fmt.Fprintln(w, strings.Join(output, "\n"))
	return err
}

func (r *ctlRunScheduleResult) RenderJSON(w io.Writer) error {
	_, err := fmt.Fprintln(w, protojson.Format(r.res))
	return err
}

var ctlRunScheduleCmd = &cobra.Command{
	Use:   "schedule <test id>",
	Short: "Schedule a test run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutputFormat(); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ctl.timeout"))
		defer cancel()

		return withControlClient(ctx, viper.GetString("ctl.grpc-addr"), !viper.GetBool("ctl.insecure"), viper.GetString("ctl.access-token"), func(ctx context.Context, client controlv1.ControlServiceClient) error {
			ctx = metadata.AppendToOutgoingContext(ctx, auth.MDKeyNamespace, ctlNamespace)
			res, err := client.ScheduleRun(ctx, &controlv1.ScheduleRunRequest{
				TestId: args[0],
			})
			if err != nil {
				return err
			}

			return render(&ctlRunScheduleResult{testID: args[0], res: res}, cmd.OutOrStdout())
		})
	},
}

func init() {
	addOutputFlag(ctlRunScheduleCmd)
	ctlRunCmd.AddCommand(ctlRunScheduleCmd)
}
