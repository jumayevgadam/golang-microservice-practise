package connection

import (
	"cart/internal/config"
	"context"
	"fmt"

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
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.PostgresConfig) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.GenerateDSN())
	if err != nil {
		return nil, fmt.Errorf("error creating pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Database{pool: pool}, nil
}

func (d *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	result, err := d.pool.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("executing query error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgconn.CommandTag{}, pgx.ErrNoRows
	}

	return result, nil
}

func (d *Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, d.pool, dest, query, args...)
}

func (d *Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, d.pool, dest, query, args...)
}

func (d *Database) Close() {
	d.pool.Close()
}
