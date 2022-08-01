package scheduler

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_generateLabelSet(t *testing.T) {
	tests := []struct {
		name         string
		baseLabels   map[string]string
		matrixLabels map[string][]string
		labelSet     []map[string]string
	}{
		{
			name:       "base, no matrix labels",
			baseLabels: map[string]string{"label": "value"},
			labelSet:   []map[string]string{{"label": "value"}},
		},
		{
			name: "no base, matrix labels",
			matrixLabels: map[string][]string{
				"a": {"1", "2", "3"},
				"b": {"4", "5"},
				"c": {"6"},
			},
			labelSet: []map[string]string{
				{"a": "1", "b": "4", "c": "6"},
				{"a": "2", "b": "4", "c": "6"},
				{"a": "3", "b": "4", "c": "6"},
				{"a": "1", "b": "5", "c": "6"},
				{"a": "2", "b": "5", "c": "6"},
				{"a": "3", "b": "5", "c": "6"},
			},
		},
		{
			name:         "matrix overwrite base",
			baseLabels:   map[string]string{"key": "value"},
			matrixLabels: map[string][]string{"key": {"value_1", "value_2"}},
			labelSet: []map[string]string{
				{"key": "value_1"},
				{"key": "value_2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labelSet := generateMatrixLabelSet(tt.baseLabels, tt.matrixLabels)

			trans := cmp.Transformer("Sort", func(in []int) []int {
				out := append([]int(nil), in...) // Copy input to avoid mutating it
				sort.Slice(out, func(i, j int) bool {
					return fmt.Sprintf("%v", out[i]) < fmt.Sprintf("%v", out[j])
				})
				return out
			})

			if !cmp.Equal(labelSet, tt.labelSet, trans) {
				t.Errorf("incorrect label set\nexpected:\t%v\ngot:\t%v", tt.labelSet, labelSet)
			}
		})
	}
}
