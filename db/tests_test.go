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
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTests(t *testing.T) {
	ctx := context.Background()

	withTestDB(t, func(db DBTX) {
		querier := New()

		var (
			testA Test
			testB Test
			testC Test
		)

		t.Run("RegisterTest", func(t *testing.T) {
			var (
				runConfig pgtype.JSONB
				labels    pgtype.JSONB
				matrix    pgtype.JSONB
			)

			err := runConfig.Set(&TestRunConfig{
				ContainerImage: "image",
				Command:        "cmd",
				Args:           []string{"a"},
				Env:            map[string]string{"env": "value"},
				TimeoutSeconds: 60,
			})
			require.NoError(t, err)

			err = labels.Set(map[string]string{"label": "value"})
			require.NoError(t, err)

			err = matrix.Set(&TestMatrix{
				Labels: map[string][]string{"label": {"value_1", "value_2"}},
			})
			require.NoError(t, err)

			testA, err = querier.RegisterTest(ctx, db, RegisterTestParams{
				Namespace:    "ns1",
				Name:         "testA",
				RunConfig:    runConfig,
				Labels:       labels,
				Matrix:       matrix,
				CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
			})
			require.NoError(t, err)

			testB, err = querier.RegisterTest(ctx, db, RegisterTestParams{
				Namespace:    "ns2",
				Name:         "testB",
				RunConfig:    runConfig,
				Labels:       labels,
				Matrix:       matrix,
				CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
			})
			require.NoError(t, err)

			testC, err = querier.RegisterTest(ctx, db, RegisterTestParams{
				Namespace:    "ns1",
				Name:         "testC",
				RunConfig:    runConfig,
				Labels:       labels,
				Matrix:       matrix,
				CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
			})
			require.NoError(t, err)
		})

		t.Run("GetTest", func(t *testing.T) {
			test, err := querier.GetTest(ctx, db, GetTestParams{
				ID:        testA.ID,
				Namespace: testA.Namespace,
			})
			require.NoError(t, err)
			assert.Equal(t, testA, test)

			_, err = querier.GetTest(ctx, db, GetTestParams{
				ID:        testB.ID,
				Namespace: "ns1",
			})
			require.Error(t, err)
			assert.ErrorIs(t, err, pgx.ErrNoRows)

			_, err = querier.GetTest(ctx, db, GetTestParams{
				ID:        testA.ID,
				Namespace: "invalid",
			})
			require.Error(t, err)
			assert.ErrorIs(t, err, pgx.ErrNoRows)
		})

		t.Run("QueryTests", func(t *testing.T) {
			t.Run("by ID", func(t *testing.T) {
				tests, err := querier.QueryTests(ctx, db, QueryTestsParams{
					Namespace: "ns1",
					Ids:       []uuid.UUID{testA.ID, testC.ID},
					Labels: pgtype.JSONB{
						Status: pgtype.Null,
					},
				})
				require.NoError(t, err)
				expected := []Test{testA, testC}
				sort.Slice(expected, func(i, j int) bool {
					return expected[i].Name < expected[j].Name
				})
				assert.Equal(t, expected, tests)

				tests, err = querier.QueryTests(ctx, db, QueryTestsParams{
					Namespace: "ns1",
					Ids:       []uuid.UUID{testB.ID},
					Labels: pgtype.JSONB{
						Status: pgtype.Null,
					},
				})
				require.NoError(t, err)
				assert.Len(t, tests, 0)
			})

			t.Run("by labels", func(t *testing.T) {
				tests, err := querier.QueryTests(ctx, db, QueryTestsParams{
					Namespace: "ns1",
					Labels: pgtype.JSONB{
						Bytes:  []byte(`{"label":"value"}`),
						Status: pgtype.Present,
					},
				})
				require.NoError(t, err)
				expected := []Test{testA, testC}
				sort.Slice(expected, func(i, j int) bool {
					return expected[i].Name < expected[j].Name
				})
				assert.Equal(t, expected, tests)

				tests, err = querier.QueryTests(ctx, db, QueryTestsParams{
					Namespace: "ns1",
					Labels: pgtype.JSONB{
						Bytes:  []byte(`{"other":"value"}`),
						Status: pgtype.Present,
					},
				})
				require.NoError(t, err)
				assert.Len(t, tests, 0)
			})
		})

		t.Run("ListAllNamespaces", func(t *testing.T) {
			ns, err := querier.ListAllNamespaces(ctx, db)
			require.NoError(t, err)
			assert.Equal(t, []string{"ns1", "ns2"}, ns)
		})

		t.Run("ListTests", func(t *testing.T) {
			tests, err := querier.ListTests(ctx, db, "ns1")
			require.NoError(t, err)
			expected := []Test{testA, testC}
			sort.Slice(expected, func(i, j int) bool {
				return expected[i].Name < expected[j].Name
			})
			assert.Equal(t, expected, tests)
		})

		t.Run("ListTestsToSchedule", func(t *testing.T) {
			// Registering a test will schedule a run for the test, so we expect no
			// tests to schedule.
			tests, err := querier.ListTestsToSchedule(ctx, db)
			require.NoError(t, err)
			assert.Equal(t, []Test(nil), tests)

			err = querier.DeleteRunsForTest(ctx, db, testA.ID)
			require.NoError(t, err)
			_, err = querier.UpdateTest(ctx, db, UpdateTestParams{
				ID:           testA.ID,
				Namespace:    testA.Namespace,
				Name:         testA.Name,
				RunConfig:    testA.RunConfig,
				Labels:       testA.Labels,
				Matrix:       testA.Matrix,
				CronSchedule: testA.CronSchedule,
				NextRunAt:    sql.NullTime{Valid: true, Time: time.Now().Add(-time.Hour)},
			})
			require.NoError(t, err)

			// Deleting the pending run for testA and changing its next run at should
			// result in it being schedulable.
			tests, err = querier.ListTestsToSchedule(ctx, db)
			require.NoError(t, err)
			assert.Len(t, tests, 1)
			assert.Equal(t, testA.ID, tests[0].ID)
		})

		t.Run("UpdateTest", func(t *testing.T) {
			var (
				runConfig pgtype.JSONB
				labels    pgtype.JSONB
				matrix    pgtype.JSONB
			)

			err := runConfig.Set(&TestRunConfig{
				ContainerImage: "new-image",
				Command:        "new-cmd",
				Args:           []string{"b"},
				Env:            map[string]string{"key": "value"},
				TimeoutSeconds: 10,
			})
			require.NoError(t, err)

			err = labels.Set(map[string]string{"one": "two"})
			require.NoError(t, err)

			err = matrix.Set(&TestMatrix{
				Labels: map[string][]string{"one": {"2", "three"}},
			})
			require.NoError(t, err)

			res, err := querier.UpdateTest(ctx, db, UpdateTestParams{
				ID:           testA.ID,
				Namespace:    "invalid",
				Name:         "new-test-a",
				RunConfig:    runConfig,
				Labels:       labels,
				Matrix:       matrix,
				CronSchedule: sql.NullString{Valid: false},
			})
			require.NoError(t, err)
			assert.Equal(t, int64(0), res.RowsAffected())

			res, err = querier.UpdateTest(ctx, db, UpdateTestParams{
				ID:           testA.ID,
				Namespace:    testA.Namespace,
				Name:         "new-test-a",
				RunConfig:    runConfig,
				Labels:       labels,
				Matrix:       matrix,
				CronSchedule: sql.NullString{Valid: false},
			})
			require.NoError(t, err)
			assert.Equal(t, int64(1), res.RowsAffected())

			test, err := querier.GetTest(ctx, db, GetTestParams{
				ID:        testA.ID,
				Namespace: testA.Namespace,
			})
			require.NoError(t, err)
			assert.Equal(t, "new-test-a", test.Name)
			assert.Equal(t, sql.NullString{}, test.CronSchedule)

			{
				var expected, got TestRunConfig
				err = runConfig.AssignTo(&expected)
				require.NoError(t, err)
				err = test.RunConfig.AssignTo(&got)
				require.NoError(t, err)
				assert.Equal(t, expected, got)
			}
			{
				var expected, got map[string]string
				err = labels.AssignTo(&expected)
				require.NoError(t, err)
				err = test.Labels.AssignTo(&got)
				require.NoError(t, err)
				assert.Equal(t, expected, got)
			}
			{
				var expected, got TestMatrix
				err = matrix.AssignTo(&expected)
				require.NoError(t, err)
				err = test.Matrix.AssignTo(&got)
				require.NoError(t, err)
				assert.Equal(t, expected, got)
			}
		})

		t.Run("DeleteTest", func(t *testing.T) {
			_, err := querier.DeleteTest(ctx, db, DeleteTestParams{
				ID:        testB.ID,
				Namespace: testB.Namespace,
			})
			require.NoError(t, err)

			_, err = querier.GetTest(ctx, db, GetTestParams{
				ID:        testB.ID,
				Namespace: testB.Namespace,
			})
			require.Error(t, err)
			assert.ErrorIs(t, err, pgx.ErrNoRows)
		})
	})
}
