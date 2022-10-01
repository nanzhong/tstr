//go:build integration

package db

import (
	"context"
	"database/sql"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunnersQueries(t *testing.T) {
	ctx := context.Background()

	withTestDB(t, func(db DBTX) {
		querier := New()

		var (
			runnerA            Runner
			runnerB            Runner
			updatedHeartbeatTS time.Time = time.Now().Add(time.Hour).Truncate(time.Second)
		)

		t.Run("RegisterRunner", func(t *testing.T) {
			var (
				acceptSelectors pgtype.JSONB
				rejectSelectors pgtype.JSONB
			)
			err := acceptSelectors.Set(map[string]string{"key": "acceptSel"})
			require.NoError(t, err)
			err = rejectSelectors.Set(map[string]string{"key": "rejectSel"})
			require.NoError(t, err)

			runnerA, err = querier.RegisterRunner(ctx, db, RegisterRunnerParams{
				Name:                     "a",
				NamespaceSelectors:       []string{"a"},
				AcceptTestLabelSelectors: acceptSelectors,
				RejectTestLabelSelectors: rejectSelectors,
			})
			require.NoError(t, err)

			runnerB, err = querier.RegisterRunner(ctx, db, RegisterRunnerParams{
				Name:                     "b",
				NamespaceSelectors:       []string{"b"},
				AcceptTestLabelSelectors: acceptSelectors,
				RejectTestLabelSelectors: rejectSelectors,
			})
			require.NoError(t, err)
		})

		t.Run("GetRunner", func(t *testing.T) {
			runner, err := querier.GetRunner(ctx, db, runnerA.ID)
			require.NoError(t, err)
			assert.Equal(t, runnerA, runner)
		})

		t.Run("UpdateRunnerHeartbeat", func(t *testing.T) {
			err := querier.UpdateRunnerHeartbeat(ctx, db, UpdateRunnerHeartbeatParams{
				ID:        runnerB.ID,
				Timestamp: sql.NullTime{Valid: true, Time: updatedHeartbeatTS},
			})
			require.NoError(t, err)

			runner, err := querier.GetRunner(ctx, db, runnerB.ID)
			require.NoError(t, err)
			assert.Equal(t, sql.NullTime{Valid: true, Time: updatedHeartbeatTS}, runner.LastHeartbeatAt)
			runnerB = runner
		})

		t.Run("ListRunners", func(t *testing.T) {
			runners, err := querier.ListRunners(ctx, db)
			require.NoError(t, err)
			expectedRunners := []Runner{runnerB, runnerA}
			sort.Slice(expectedRunners, func(i, j int) bool {
				return expectedRunners[i].LastHeartbeatAt.Time.After(expectedRunners[j].LastHeartbeatAt.Time)
			})
			assert.Equal(t, expectedRunners, runners)
		})

		t.Run("QueryRunners", func(t *testing.T) {
			t.Run("by ID", func(t *testing.T) {
				runners, err := querier.QueryRunners(ctx, db, QueryRunnersParams{
					Ids: []uuid.UUID{runnerA.ID},
				})
				require.NoError(t, err)
				assert.Equal(t, []Runner{runnerA}, runners)
			})

			t.Run("by last heartbeat", func(t *testing.T) {
				runners, err := querier.QueryRunners(ctx, db, QueryRunnersParams{
					LastHeartbeatSince: sql.NullTime{Valid: true, Time: updatedHeartbeatTS.Add(-time.Second)},
				})
				require.NoError(t, err)
				assert.Equal(t, []Runner{runnerB}, runners)
			})
		})
	})
}
