// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

const getRunnerSQL = `SELECT *
FROM runners
WHERE id = $1;`

type GetRunnerRow struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	AcceptTestLabels pgtype.JSONB `json:"accept_test_labels"`
	RejectTestLabels pgtype.JSONB `json:"reject_test_labels"`
	RegisteredAt     time.Time    `json:"registered_at"`
	ApprovedAt       time.Time    `json:"approved_at"`
	RevokedAt        time.Time    `json:"revoked_at"`
	LastHeartbeatAt  time.Time    `json:"last_heartbeat_at"`
}

// GetRunner implements Querier.GetRunner.
func (q *DBQuerier) GetRunner(ctx context.Context, id string) (GetRunnerRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetRunner")
	row := q.conn.QueryRow(ctx, getRunnerSQL, id)
	var item GetRunnerRow
	if err := row.Scan(&item.ID, &item.Name, &item.AcceptTestLabels, &item.RejectTestLabels, &item.RegisteredAt, &item.ApprovedAt, &item.RevokedAt, &item.LastHeartbeatAt); err != nil {
		return item, fmt.Errorf("query GetRunner: %w", err)
	}
	return item, nil
}

// GetRunnerBatch implements Querier.GetRunnerBatch.
func (q *DBQuerier) GetRunnerBatch(batch genericBatch, id string) {
	batch.Queue(getRunnerSQL, id)
}

// GetRunnerScan implements Querier.GetRunnerScan.
func (q *DBQuerier) GetRunnerScan(results pgx.BatchResults) (GetRunnerRow, error) {
	row := results.QueryRow()
	var item GetRunnerRow
	if err := row.Scan(&item.ID, &item.Name, &item.AcceptTestLabels, &item.RejectTestLabels, &item.RegisteredAt, &item.ApprovedAt, &item.RevokedAt, &item.LastHeartbeatAt); err != nil {
		return item, fmt.Errorf("scan GetRunnerBatch row: %w", err)
	}
	return item, nil
}

const listRunnersSQL = `SELECT *
FROM runners
WHERE
  CASE WHEN $1
    THEN revoked_at IS NOT NULL
    ELSE TRUE
  END;`

type ListRunnersRow struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	AcceptTestLabels pgtype.JSONB `json:"accept_test_labels"`
	RejectTestLabels pgtype.JSONB `json:"reject_test_labels"`
	RegisteredAt     time.Time    `json:"registered_at"`
	ApprovedAt       time.Time    `json:"approved_at"`
	RevokedAt        time.Time    `json:"revoked_at"`
	LastHeartbeatAt  time.Time    `json:"last_heartbeat_at"`
}

// ListRunners implements Querier.ListRunners.
func (q *DBQuerier) ListRunners(ctx context.Context, filterRevoked bool) ([]ListRunnersRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ListRunners")
	rows, err := q.conn.Query(ctx, listRunnersSQL, filterRevoked)
	if err != nil {
		return nil, fmt.Errorf("query ListRunners: %w", err)
	}
	defer rows.Close()
	items := []ListRunnersRow{}
	for rows.Next() {
		var item ListRunnersRow
		if err := rows.Scan(&item.ID, &item.Name, &item.AcceptTestLabels, &item.RejectTestLabels, &item.RegisteredAt, &item.ApprovedAt, &item.RevokedAt, &item.LastHeartbeatAt); err != nil {
			return nil, fmt.Errorf("scan ListRunners row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListRunners rows: %w", err)
	}
	return items, err
}

// ListRunnersBatch implements Querier.ListRunnersBatch.
func (q *DBQuerier) ListRunnersBatch(batch genericBatch, filterRevoked bool) {
	batch.Queue(listRunnersSQL, filterRevoked)
}

// ListRunnersScan implements Querier.ListRunnersScan.
func (q *DBQuerier) ListRunnersScan(results pgx.BatchResults) ([]ListRunnersRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query ListRunnersBatch: %w", err)
	}
	defer rows.Close()
	items := []ListRunnersRow{}
	for rows.Next() {
		var item ListRunnersRow
		if err := rows.Scan(&item.ID, &item.Name, &item.AcceptTestLabels, &item.RejectTestLabels, &item.RegisteredAt, &item.ApprovedAt, &item.RevokedAt, &item.LastHeartbeatAt); err != nil {
			return nil, fmt.Errorf("scan ListRunnersBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListRunnersBatch rows: %w", err)
	}
	return items, err
}

const approveRunnerSQL = `UPDATE runners
SET approved_at = CURRENT_TIMESTAMP
WHERE id = $1::uuid;`

// ApproveRunner implements Querier.ApproveRunner.
func (q *DBQuerier) ApproveRunner(ctx context.Context, id string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ApproveRunner")
	cmdTag, err := q.conn.Exec(ctx, approveRunnerSQL, id)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query ApproveRunner: %w", err)
	}
	return cmdTag, err
}

// ApproveRunnerBatch implements Querier.ApproveRunnerBatch.
func (q *DBQuerier) ApproveRunnerBatch(batch genericBatch, id string) {
	batch.Queue(approveRunnerSQL, id)
}

// ApproveRunnerScan implements Querier.ApproveRunnerScan.
func (q *DBQuerier) ApproveRunnerScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec ApproveRunnerBatch: %w", err)
	}
	return cmdTag, err
}

const revokeRunnerSQL = `UPDATE runners
SET revoked_at = CURRENT_TIMESTAMP
WHERE id = $1::uuid;`

// RevokeRunner implements Querier.RevokeRunner.
func (q *DBQuerier) RevokeRunner(ctx context.Context, id string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "RevokeRunner")
	cmdTag, err := q.conn.Exec(ctx, revokeRunnerSQL, id)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query RevokeRunner: %w", err)
	}
	return cmdTag, err
}

// RevokeRunnerBatch implements Querier.RevokeRunnerBatch.
func (q *DBQuerier) RevokeRunnerBatch(batch genericBatch, id string) {
	batch.Queue(revokeRunnerSQL, id)
}

// RevokeRunnerScan implements Querier.RevokeRunnerScan.
func (q *DBQuerier) RevokeRunnerScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec RevokeRunnerBatch: %w", err)
	}
	return cmdTag, err
}
