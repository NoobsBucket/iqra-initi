package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(dsn string, maxConnection int, maxIdelConnection int, maxIdleTime time.Duration) (*pgxpool.Pool, error) {

	cfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	cfg.MaxConnIdleTime = maxIdleTime
	cfg.MaxConns = int32(maxConnection)
	cfg.MinConns = int32(maxIdelConnection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil

}
