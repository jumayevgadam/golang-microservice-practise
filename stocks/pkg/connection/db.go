package connection

import (
	"context"
	"fmt"
	"stocks/internal/config"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Ensure Database implements DB and DBOps.
var _ DB = (*Database)(nil)

// Querier defines basic query operations.
type Querier interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

// DB extends Querier with additional convenience methods.
type DB interface {
	Querier
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Close()
}

type Database struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.PostgresConfig) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.GenerateDSN())
	if err != nil {
		return nil, fmt.Errorf("error creating pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

func (d *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return d.Pool.QueryRow(ctx, query, args...)
}

func (d *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.Pool.Query(ctx, query, args...)
}

func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	result, err := d.Pool.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("executing query error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgconn.CommandTag{}, pgx.ErrNoRows
	}

	return result, nil
}

func (d *Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, d.Pool, dest, query, args...)
}

func (d *Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, d.Pool, dest, query, args...)
}

func (d *Database) Close() {
	d.Pool.Close()
}
