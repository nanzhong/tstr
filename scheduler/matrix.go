package scheduler

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/db"
)

func RunsForTest(test db.Test) ([]db.ScheduleRunParams, error) {
	var labels map[string]string
	if err := test.Labels.AssignTo(&labels); err != nil {
		return nil, fmt.Errorf("parsing labels: %w", err)
	}

	var matrix db.TestMatrix
	if err := test.Matrix.AssignTo(&matrix); err != nil {
		return nil, fmt.Errorf("parsing matrix: %w", err)
	}

	matrixID := uuid.NullUUID{}
	labelSet := generateMatrixLabelSet(labels, matrix.Labels)
	if len(matrix.Labels) > 0 {
		matrixID.Valid = true
		matrixID.UUID = uuid.New()
	}

	var runParams []db.ScheduleRunParams
	for _, labels := range labelSet {
		var dbLabels pgtype.JSONB
		if err := dbLabels.Set(labels); err != nil {
			return nil, fmt.Errorf("formatting labels: %w", err)
		}
		runParams = append(runParams, db.ScheduleRunParams{
			Labels:       dbLabels,
			TestMatrixID: matrixID,
			TestID:       test.ID,
		})
	}
	return runParams, nil
}

func generateMatrixLabelSet(b map[string]string, ml map[string][]string) []map[string]string {
	var labelSets []map[string]string
	if len(ml) == 0 {
		labelSets = append(labelSets, copyLabels(b))
	} else {
		var (
			kIdx = make(map[string]int)
			keys []string
		)
		for k := range ml {
			kIdx[k] = 0
			keys = append(keys, k)
		}

		for {
			labels := make(map[string]string)
			for _, rk := range keys {
				labels[rk] = ml[rk][kIdx[rk]]
			}
			labelSets = append(labelSets, labels)

			next := -1
			for i, k := range keys {
				if kIdx[k]+1 < len(ml[k]) {
					next = i
					break
				}
			}
			if next == -1 {
				break
			}
			kIdx[keys[next]] += 1
			for i := 0; i < next; i++ {
				kIdx[keys[i]] = 0
			}
		}
	}

	return labelSets
}
