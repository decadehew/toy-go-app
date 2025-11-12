package db

// import (
// 	"context"
// 	"database/sql"
// 	"time"
// )

// func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", addr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(maxOpenConns)
// 	duration, err := time.ParseDuration(maxIdleTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.SetConnMaxIdleTime(duration)
// 	db.SetMaxIdleConns(maxIdleConns)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err = db.PingContext(ctx); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxOpenConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	config.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
