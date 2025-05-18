package sqlite

import (
	"context"
	"database/sql"

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
	return nil
}

func (c *client) OnStop(context.Context) error {
	return c.pool.Close()
}

func (c *client) DB() *sql.DB {
	return c.pool
}
