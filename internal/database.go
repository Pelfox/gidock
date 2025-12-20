package internal

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreatePool creates and configures a new PostgreSQL connection pool using
// `pgxpool.Pool`.
func CreatePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// setting up pool config
	// TODO: make this settings configurable
	config.MinConns = 2
	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(timedCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
