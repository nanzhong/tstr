package scheduler

import (
	"sort"
	"strings"
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
				{"a": "1", "b": "5", "c": "6"},
				{"a": "2", "b": "4", "c": "6"},
				{"a": "2", "b": "5", "c": "6"},
				{"a": "3", "b": "4", "c": "6"},
				{"a": "3", "b": "5", "c": "6"},
			},
		},
		{
			name:         "matrix overwrite base",
			baseLabels:   map[string]string{"key": "value", "other": "label"},
			matrixLabels: map[string][]string{"key": {"value_1", "value_2"}},
			labelSet: []map[string]string{
				{"key": "value_1", "other": "label"},
				{"key": "value_2", "other": "label"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labelSet := generateMatrixLabelSet(tt.baseLabels, tt.matrixLabels)

			trans := cmp.Transformer("SortLabels", func(labels []map[string]string) []string {
				var lStrs []string
				for _, l := range labels {
					ls, err := sortedLabelsString(l)
					if err != nil {
						t.Fatal(err)
					}
					lStrs = append(lStrs, ls)
				}
				sort.Strings(lStrs)
				return lStrs
			})
			if !cmp.Equal(labelSet, tt.labelSet, trans) {
				t.Errorf("incorrect label set\nexpected:\t%v\ngot:\t%v", tt.labelSet, labelSet)
			}
		})
	}
}

func sortedLabelsString(l map[string]string) (string, error) {
	var keys []string
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		_, err := sb.WriteString("[" + k + "=" + l[k] + "]")
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}
