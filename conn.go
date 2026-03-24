package qube

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
)

type Conn struct {
	db  *sql.DB
	raw *sql.Conn
}

func (conn *Conn) Exec(query string, args ...any) (sql.Result, error) {
	// Avoid "bad connection".
	return conn.withRetry(context.Background(), query, args...)
}

func (conn *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	// Avoid "bad connection".
	return conn.withRetry(ctx, query, args...)
}

func (conn *Conn) withRetry(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := conn.raw.ExecContext(ctx, query, args...)

	if errors.Is(err, driver.ErrBadConn) {
		conn.raw.Close()
		raw, err := conn.db.Conn(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to reopen DB connection (%w)", err)
		}

		conn.raw = raw
		return conn.raw.ExecContext(ctx, query, args...)
	}

	return res, err
}

func (conn *Conn) Close() {
	conn.raw.Close()
	conn.db.Close()
}
