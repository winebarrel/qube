package qube

import (
	"context"
	"database/sql"
)

type DBIface interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Close() error
}
