package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/shamil/weather/internal/infrastructure/database"
)

type Pool struct {
	db *sql.DB
}

func (c *Pool) Builder() *sql.DB {
	return c.db
}

// Drop close not implemented in database
func (c *Pool) Drop() error {
	return nil
}

func (c *Pool) DropMsg() string {
	return "close database: is not implemented"
}

func NewPool(ctx context.Context, opt *database.Opt) (*Pool, error) {
	db, err := sql.Open(opt.Dialect, opt.ConnectionString())
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = db.PingContext(pingCtx); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opt.MaxIdleConns)
	db.SetMaxOpenConns(opt.MaxOpenConns)
	db.SetConnMaxLifetime(opt.MaxConnMaxLifetime)

	return &Pool{db: db}, nil
}
