package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pgxPool *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := initDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %s", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func initDB() error {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return errors.New("missing TEST_DATABASE_URL env var")
	}

	var err error
	pgxPool, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	return nil
}

func withTestDB(tb testing.TB, fn func(db DBTX)) {
	tx, err := pgxPool.Begin(context.Background())
	if err != nil {
		tb.Fatalf("failed to create test transaction: %s", err)
	}
	if err := tx.Rollback(context.Background()); err != nil {
		tb.Fatalf("failed to rollback test transaction: %s", err)
	}
}
