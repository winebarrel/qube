package qube

import (
	"context"
	"database/sql"
	"fmt"
	"io"
)

type NullDB struct {
	w io.Writer
}

func (db *NullDB) Exec(query string, args ...any) (sql.Result, error) {
	fmt.Fprintln(db.w, query)
	return nil, nil
}

func (db *NullDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	fmt.Fprintln(db.w, query)
	return nil, nil
}

func (db *NullDB) Close() error {
	return nil
}
