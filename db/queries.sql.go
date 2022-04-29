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

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	RegisterTest(ctx context.Context, params RegisterTestParams) (RegisterTestRow, error)
	// RegisterTestBatch enqueues a RegisterTest query into batch to be executed
	// later by the batch.
	RegisterTestBatch(batch genericBatch, params RegisterTestParams)
	// RegisterTestScan scans the result of an executed RegisterTestBatch query.
	RegisterTestScan(results pgx.BatchResults) (RegisterTestRow, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...interface{})
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, registerTestSQL, registerTestSQL); err != nil {
		return fmt.Errorf("prepare query 'RegisterTest': %w", err)
	}
	return nil
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

const registerTestSQL = `WITH data (name, labels, cron_schedule, container_image, command, args, env) AS (
  VALUES (
    $1::varchar,
    $2::jsonb,
    $3::varchar,
    $4::varchar,
    $5::varchar,
    $6::varchar[],
    $7::jsonb
  )
), test AS (
  INSERT INTO tests (name, labels, cron_schedule)
  SELECT name, labels, cron_schedule
  FROM data
  RETURNING id, name, labels, cron_schedule, registered_at, updated_at
), test_run_config AS (
  INSERT INTO test_run_configs (test_id, container_image, command, args, env)
  SELECT test.id, container_image, command, args, env
  FROM data, test
  RETURNING id AS test_run_config_version, container_image, command, args, env, created_at AS test_run_config_created_at
)
SELECT * FROM test, test_run_config;`

type RegisterTestParams struct {
	Name           string
	Labels         pgtype.JSONB
	CronSchedule   string
	ContainerImage string
	Command        string
	Args           []string
	Env            pgtype.JSONB
}

type RegisterTestRow struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Labels                 pgtype.JSONB `json:"labels"`
	CronSchedule           string       `json:"cron_schedule"`
	RegisteredAt           time.Time    `json:"registered_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
	TestRunConfigVersion   int          `json:"test_run_config_version"`
	ContainerImage         string       `json:"container_image"`
	Command                string       `json:"command"`
	Args                   []string     `json:"args"`
	Env                    pgtype.JSONB `json:"env"`
	TestRunConfigCreatedAt time.Time    `json:"test_run_config_created_at"`
}

// RegisterTest implements Querier.RegisterTest.
func (q *DBQuerier) RegisterTest(ctx context.Context, params RegisterTestParams) (RegisterTestRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "RegisterTest")
	row := q.conn.QueryRow(ctx, registerTestSQL, params.Name, params.Labels, params.CronSchedule, params.ContainerImage, params.Command, params.Args, params.Env)
	var item RegisterTestRow
	if err := row.Scan(&item.ID, &item.Name, &item.Labels, &item.CronSchedule, &item.RegisteredAt, &item.UpdatedAt, &item.TestRunConfigVersion, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.TestRunConfigCreatedAt); err != nil {
		return item, fmt.Errorf("query RegisterTest: %w", err)
	}
	return item, nil
}

// RegisterTestBatch implements Querier.RegisterTestBatch.
func (q *DBQuerier) RegisterTestBatch(batch genericBatch, params RegisterTestParams) {
	batch.Queue(registerTestSQL, params.Name, params.Labels, params.CronSchedule, params.ContainerImage, params.Command, params.Args, params.Env)
}

// RegisterTestScan implements Querier.RegisterTestScan.
func (q *DBQuerier) RegisterTestScan(results pgx.BatchResults) (RegisterTestRow, error) {
	row := results.QueryRow()
	var item RegisterTestRow
	if err := row.Scan(&item.ID, &item.Name, &item.Labels, &item.CronSchedule, &item.RegisteredAt, &item.UpdatedAt, &item.TestRunConfigVersion, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.TestRunConfigCreatedAt); err != nil {
		return item, fmt.Errorf("scan RegisterTestBatch row: %w", err)
	}
	return item, nil
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0