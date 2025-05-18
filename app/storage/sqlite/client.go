package sqlite

import (
	"context"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type client struct {
	dsn  string
	pool *sql.DB
}

func NewClient(dsn string) *client {
	return &client{
		dsn:  dsn,
		pool: nil,
	}
}

func (c *client) OnStart(ctx context.Context) error {
	db, err := sql.Open("sqlite", c.dsn)
	if err != nil {
		return err
	}
	c.pool = db

	return c.ping()
}

func (c *client) OnStop(context.Context) error {
	return c.pool.Close()
}

func (c *client) DB() *sql.DB {
	return c.pool
}

func (c *client) ping() error {
	var v int
	if err := c.pool.QueryRow("SELECT 1").Scan(&v); err != nil {
		return err
	} else if v != 1 {
		return errors.New("sqlite.client: health check failed")
	}
	return nil
}
