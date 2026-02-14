package qube

import (
	"context"
	"database/sql"
)

type DBIface interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Close() error
}
